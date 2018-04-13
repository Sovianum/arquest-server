package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadConf(r io.Reader) (*Conf, error) {
	conf := &Conf{}
	parseErr := json.NewDecoder(r).Decode(&conf)
	return conf, parseErr
}

type Conf struct {
	Log         string      `json:"log"`
	PortEnvVar  string      `json:"port_env_var"`
	DefaultPort int         `json:"default_port"`
	Auth        AuthConfig  `json:"auth"`
	DB          DBConfig    `json:"db"`
	Logic       LogicConfig `json:"logic"`
}

type AuthConfig struct {
	TokenKey   string `json:"token_key"`
	ExpireDays int    `json:"expire_days"`
}

type DBConfig struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	EnvVar             string `json:"env_var"`
	DriverName         string `json:"driver_name"`
	User               string `json:"user"`
	Password           string `json:"password"`
	DBName             string `json:"db_name"`
	AuthStringTemplate string `json:"auth_string_template"`
}

type LogicConfig struct {
	QuestDataTemplate string `json:"quest_data_template"`
}

func (conf AuthConfig) GetTokenKey() []byte {
	return []byte(conf.TokenKey) // TODO use secure service instead of bicycles
}

func (conf DBConfig) GetAuthStr() string {
	return fmt.Sprintf(conf.AuthStringTemplate, conf.Host, conf.Port, conf.User, conf.Password, conf.DBName)
}

func (conf DBConfig) GetEnvAuthString() string {
	return os.Getenv(conf.EnvVar)
}
