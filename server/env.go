package server

import (
	"database/sql"
	"github.com/Sovianum/acquaintanceServer/config"
	"github.com/Sovianum/acquaintanceServer/dao"
)

type tokenKeyGetterType func() string

type Env struct {
	userDAO       dao.UserDAO
	positionDAO   dao.PositionDAO
	authConf      config.AuthConfig
	hashFunc      func(password []byte) ([]byte, error)
	hashValidator func(password []byte, hash []byte) error
}

func NewEnv(db *sql.DB, authConf config.AuthConfig) *Env {
	var result = &Env{
		userDAO:     dao.NewDBUserDAO(db),
		positionDAO: dao.NewDBPositionDAO(db),
		authConf:    authConf,
	}

	return result
}
