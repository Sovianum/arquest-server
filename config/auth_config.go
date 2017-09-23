package config

import (
	"encoding/json"
	"io"
)

func ReadAuthConf(r io.Reader) (AuthConfig, error) {
	var conf = AuthConfig{}

	var parseErr = json.NewDecoder(r).Decode(&conf)
	return conf, parseErr
}

type AuthConfig struct {
	TokenKey   string `json:"token_key"`
	ExpireDays int    `json:"expire_days"`
}

func (conf *AuthConfig) GetTokenKey() []byte {
	return []byte(conf.TokenKey) // TODO use secure service instead of bicycles
}
