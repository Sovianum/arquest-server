package common

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

type ResponseMsg struct {
	ErrMsg interface{} `json:"err_msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func GetErrResponse(err error) ResponseMsg {
	return ResponseMsg{ErrMsg: err.Error()}
}

func GetDataResponse(data interface{}) ResponseMsg {
	return ResponseMsg{Data: data}
}

func GetEmptyResponse() ResponseMsg {
	return ResponseMsg{}
}
