package config

import (
	"encoding/json"
	"fmt"
	"io"
)

func ReadDBConfig(r io.Reader) (DBConfig, error) {
	var conf = DBConfig{}

	var parseErr = json.NewDecoder(r).Decode(&conf)
	return conf, parseErr
}

type DBConfig struct {
	Port               int    `json:"port"`
	DriverName         string `json:"driver_name"`
	User               string `json:"user"`
	Password           string `json:"password"`
	DBName             string `json:"db_name"`
	AuthStringTemplate string `json:"auth_string_template"`
}

func (conf DBConfig) GetAuthStr() string {
	return fmt.Sprintf(conf.AuthStringTemplate, conf.User, conf.Password, conf.DBName)
}
