package config

import (
	"io"
	"encoding/json"
	"fmt"
)

func ReadConf(r io.Reader) (Conf, error) {
	var conf = Conf{}

	var parseErr = json.NewDecoder(r).Decode(&conf)
	return conf, parseErr
}

type Conf struct {
	Auth  AuthConfig  `json:"auth"`
	DB    DBConfig    `json:"db"`
	Logic LogicConfig `json:"logic"`
}

type AuthConfig struct {
	TokenKey   string `json:"token_key"`
	ExpireDays int    `json:"expire_days"`
}

type DBConfig struct {
	Port               int    `json:"port"`
	DriverName         string `json:"driver_name"`
	User               string `json:"user"`
	Password           string `json:"password"`
	DBName             string `json:"db_name"`
	AuthStringTemplate string `json:"auth_string_template"`
}

type LogicConfig struct {
	Distance          float64 `json:"distance"`
	OnlineTimeout     int     `json:"online_timeout"`
	RequestExpiration int     `json:"request_expiration"`
	CleanupInterval   int     `json:"cleanup_interval"`
}

func (conf AuthConfig) GetTokenKey() []byte {
	return []byte(conf.TokenKey) // TODO use secure service instead of bicycles
}

func (conf DBConfig) GetAuthStr() string {
	return fmt.Sprintf(conf.AuthStringTemplate, conf.User, conf.Password, conf.DBName)
}
