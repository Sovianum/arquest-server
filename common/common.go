package common

import "encoding/json"

type ResponseMsg struct {
	ErrMsg interface{} `json:"err_msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func GetErrorJson(err error) []byte {
	var msg, _ = json.Marshal(ResponseMsg{ErrMsg: err.Error()})
	return msg
}

func GetDataJson(data interface{}) []byte {
	var msg, _ = json.Marshal(ResponseMsg{Data: data})
	return msg
}

func GetEmptyJson() []byte {
	return []byte("{}")
}
