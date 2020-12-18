package gows

import (
	"encoding/json"
)

type ActionForm struct {
	Action string `json:"action"`
}

type WsForm struct {
	ActionForm
	Data interface{}
}

func GetForm(data []byte,stu interface{}) (err error) {
	err=json.Unmarshal(data,&stu)
	return
}





