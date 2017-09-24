package common

import "encoding/json"

func GetErrorJson(err error) []byte {
	var resStruct = struct {
		ErrMsg string
	}{
		ErrMsg: err.Error(),
	}

	var msg, _ = json.Marshal(resStruct)
	return msg
}
