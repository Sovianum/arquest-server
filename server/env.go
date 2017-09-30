package server

import (
	"github.com/Sovianum/acquaintanceServer/config"
	"github.com/Sovianum/acquaintanceServer/dao"
	"github.com/patrickmn/go-cache"
)

type tokenKeyGetterType func() string

type Env struct {
	userDAO          dao.UserDAO
	positionDAO      dao.PositionDAO
	meetRequestDAO   dao.MeetRequestDAO
	conf             config.Conf
	hashFunc         func(password []byte) ([]byte, error)
	hashValidator    func(password []byte, hash []byte) error
	meetRequestCache cache.Cache
}
