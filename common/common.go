package common

import (
	"encoding/json"
	"github.com/Sovianum/arquest-server/mylog"
	"net/http"
)

type ResponseMsg struct {
	ErrMsg interface{} `json:"err_msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func Round(f float64) int {
	floor := int(f)
	ceil := floor + 1

	floorDiff := f - float64(floor)
	ceilDiff := float64(ceil) - f

	result := ceil
	if floorDiff < ceilDiff {
		result = floor
	}
	return result
}

func GetErrorJson(err error) []byte {
	msg, _ := json.Marshal(ResponseMsg{ErrMsg: err.Error()})
	return msg
}

func GetDataJson(data interface{}) []byte {
	msg, _ := json.Marshal(ResponseMsg{Data: data})
	return msg
}

func GetEmptyJson() []byte {
	return []byte("{}")
}

func WriteWithLogging(r *http.Request, w http.ResponseWriter, body []byte, logger *mylog.Logger) {
	logger.LogResponseBody(r, string(body))
	w.Write(body)
}
