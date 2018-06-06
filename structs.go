package main

type UpdateDevice struct {
	Command string `json:"command"`
	Idx     int    `json:"idx"`
	Nvalue  int    `json:"nvalue"`
	Svalue  string `json:"svalue"`
}

/*func (u *UpdateDevice) MarshalJson() ([]byte, error){
    type Alias UpdateDevice
    return json.Marshal(&struct {
        Svalue string `json:"svalue"`
        *Alias
    }{
        Svalue: strings.Join(u.Svalue, ";"),
        Alias:  (*Alias)(u),
    })
}*/
