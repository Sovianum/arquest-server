package server

import (
	"github.com/Sovianum/acquaintanceServer/config"
	"github.com/Sovianum/acquaintanceServer/dao"
	"github.com/patrickmn/go-cache"
	"database/sql"
	"time"
	"crypto/sha256"
	"fmt"
)

type tokenKeyGetterType func() string

func NewEnv(db *sql.DB, conf config.Conf) *Env {
	var env = &Env{
		userDAO:dao.NewDBUserDAO(db),
		positionDAO:dao.NewDBPositionDAO(db),
		meetRequestDAO:dao.NewMeetDAO(db),
		conf:conf,
		meetRequestCache:cache.New(
			time.Second * time.Duration(conf.Logic.RequestExpiration),
			time.Second * time.Duration(conf.Logic.CleanupInterval),
		),
		hashFunc: func(password []byte) ([]byte, error) {
			var h = sha256.New()
			h.Write(password)
			return h.Sum(nil), nil
		},
		hashValidator: func(password []byte, hash []byte) error {
			var h = sha256.New()
			h.Write(password)
			var passHash = h.Sum(nil)

			if string(passHash) != string(hash) {
				return fmt.Errorf("hashes %s, %s do not match", string(passHash), string(hash))
			}
			return nil
		},
	}

	env.RunDaemons()
	return env
}

type Env struct {
	userDAO          dao.UserDAO
	positionDAO      dao.PositionDAO
	meetRequestDAO   dao.MeetRequestDAO
	conf             config.Conf
	hashFunc         func(password []byte) ([]byte, error)
	hashValidator    func(password []byte, hash []byte) error
	meetRequestCache *cache.Cache
}
