package config

import (
	"encoding/json"
	"io"
)

func ReadLogicConfig(r io.Reader) (LogicConfig, error) {
	var conf = LogicConfig{}

	var parseErr = json.NewDecoder(r).Decode(&conf)
	return conf, parseErr
}

type LogicConfig struct {
	Distance      float64 `json:"distance"`
	OnlineTimeout int     `json:"online_timeout"`
}
