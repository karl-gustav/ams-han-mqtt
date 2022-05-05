test::
	go test ./...
build::
	GOOS=linux GOARCH=arm GOARM=6 go build -o ams-han-mqtt.new
deploy:: build
	./deploy.sh
setup::
	./setup.sh
