package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/serial"
	ams "github.com/karl-gustav/ams-han"
)

var (
	address  string
	mqttUrl  string
	baudrate int
	databits int
	stopbits int
	parity   string
	verbose  bool
)

const (
	ampereSensorIdx = 610
	voltSensorIdx   = 613
	usageSensorIdx  = 615
	divider         = 10.0
)

func main() {
	flag.StringVar(&mqttUrl, "m", "tcp://localhost:1883", "url for the mqtt broker")
	flag.StringVar(&address, "a", "/dev/ttyUSB0", "Mbus serial adapter address (e.g. /dev/ttyUSB0")
	flag.IntVar(&baudrate, "b", 2400, "baud rate")
	flag.IntVar(&databits, "d", 8, "data bits")
	flag.IntVar(&stopbits, "s", 1, "stop bits")
	flag.StringVar(&parity, "p", "E", "parity (N/E/O)")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	mqtt := setupMqtt(mqttUrl)

	serialPort := getSerialPort(address, baudrate, databits, stopbits, parity)
	byteStream := createByteStream(serialPort)
	next := ams.ByteReader(byteStream)

	var usageCounter int
	for {
		bytePackage, err := next()
		if err != nil {
			fmt.Println("[ERROR]", err)
			if err == ams.CHANNEL_IS_CLOSED_ERROR {
				panic(err)
			}
			continue
		}

		if verbose {
			fmt.Printf("\nBuffer(%d): \n[%s]\n", len(bytePackage), strings.Join(byteArrayToHexStringArray(bytePackage), ","))
		}
		message, err := ams.BytesParser(bytePackage)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		switch item := message.(type) {
		case *ams.ThreeFasesMessageType2:
			mqtt <- marshalCommand(UpdateDevice{
				Command: "udevice",
				Idx:     ampereSensorIdx,
				Nvalue:  0,
				Svalue:  fmt.Sprintf("%.1f;%.1f;%.1f", float64(item.CurrL1)/divider, float64(item.CurrL2)/divider, float64(item.CurrL3)/divider),
			})
			mqtt <- marshalCommand(UpdateDevice{
				Command: "udevice",
				Idx:     voltSensorIdx,
				Nvalue:  0,
				Svalue:  fmt.Sprintf("%.3f", (float64(item.VoltL1+item.VoltL2+item.VoltL3))/30),
			})
			if usageCounter != 0 {
				mqtt <- marshalCommand(UpdateDevice{
					Command: "udevice",
					Idx:     usageSensorIdx,
					Nvalue:  0,
					Svalue:  fmt.Sprintf("%d;%d", item.ActPowPos, usageCounter),
				})
			}
		case *ams.ThreeFasesMessageType3:
			usageCounter = item.ActEnergyPa
		}

		if verbose {
			j, _ := json.Marshal(message)
			fmt.Printf("%s\n", j)
		}
	}
}

func setupMqtt(brokerUrl string) chan<- []byte {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	outgoingMessages := make(chan []byte)
	connectionLostCounter := 0
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerUrl)
	opts.SetClientID("ams-han-mqtt-client-" + strconv.Itoa(time.Now().Nanosecond()))
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		connectionLostCounter++
		message := fmt.Sprintf("MQTT connection lost (%d times): %v", connectionLostCounter, err)
		fmt.Println(message)
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Println("Connected to MQTT server!")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		for message := range outgoingMessages {
			client.Publish("domoticz/in", 2, false, message)
		}
		defer client.Disconnect(250)
	}()
	return outgoingMessages
}

func marshalCommand(command interface{}) []byte {
	commandString, err := json.Marshal(command)
	if err != nil {
		fmt.Printf("Couldn't marshal command: %+v\n", command)
	}
	return commandString
}

func getSerialPort(Address string, BaudRate int, DataBits int, StopBits int, Parity string) (port serial.Port) {
	config := serial.Config{
		Address:  address,
		BaudRate: baudrate,
		DataBits: databits,
		StopBits: stopbits,
		Parity:   parity,
		Timeout:  60 * time.Second,
	}
	if verbose {
		log.Printf("connecting %+v\n", config)
	}
	port, err := serial.Open(&config)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("connected")
	}
	return
}

func createByteStream(port serial.Port) chan byte {
	serialChannel := make(chan byte)

	go func() {
		var buf [8]byte
		for {
			n, err := port.Read(buf[:])
			if err == io.EOF {
				log.Fatalln("Reached end of stream")
				break
			} else if err != nil {
				log.Println("[ERROR]:", err)
				break
			}
			for i := 0; i < n; i++ {
				serialChannel <- buf[i]
			}
		}

		err := port.Close()
		log.Println("Closed connection!")
		close(serialChannel)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return serialChannel
}

func byteArrayToHexStringArray(bytes []byte) (strings []string) {
	for _, b := range bytes {
		strings = append(strings, fmt.Sprintf("0x%02x", b))
	}
	return
}
