package server

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/config"
	"github.com/Sovianum/arquest-server/sqldao"
	"github.com/Sovianum/arquest-server/dao"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

const (
	authorizationStr = "Authorization"
)

type tokenKeyGetterType func() string

func NewSQLEnv(db *sql.DB, conf *config.Conf, logger *mylog.Logger) *Env {
	env := &Env{
		userDAO:  sqldao.NewDBUserDAO(db),
		questDAO: sqldao.NewQuestDAO(db),
		markDAO:  sqldao.NewMarkDAO(db),
		conf:     conf,
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
		logger: logger,
	}
	return env
}

type Env struct {
	userDAO       dao.UserDAO
	questDAO      dao.QuestDAO
	markDAO       dao.MarkDAO
	conf          *config.Conf
	hashFunc      func(password []byte) ([]byte, error)
	hashValidator func(password []byte, hash []byte) error
	logger        *mylog.Logger
}

// TODO use some standard mechanisms instead of bicycles
func (env *Env) parseTokenString(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return env.conf.Auth.GetTokenKey(), nil
	})
}

func (env *Env) getIdFromRequest(r *http.Request) (id int, code int, err error) {
	headers := r.Header
	authHeaderList, ok := headers[authorizationStr]
	if !ok {
		return 0, http.StatusUnauthorized, fmt.Errorf("header \"Authorization\" not set in request")
	}
	if len(authHeaderList) != 1 {
		return 0, http.StatusBadRequest, fmt.Errorf("you set too many (%d) \"Authorization\" headers", len(authHeaderList))
	}
	authHeader := authHeaderList[0]

	fields := strings.Fields(authHeader) // getting last word to remove Bearer word from header
	tokenString := fields[len(fields)-1]

	token, tokenErr := env.parseTokenString(tokenString)
	if tokenErr != nil {
		return 0, http.StatusBadRequest, fmt.Errorf("you sent unparseable token")
	}

	userId, idErr := env.getIdFromTokenString(token)
	if idErr != nil {
		return 0, http.StatusBadRequest, fmt.Errorf("your token does not contain your id")
	}

	return userId, http.StatusOK, nil
}

func (env *Env) getIdFromTokenString(token *jwt.Token) (int, error) {
	claims, okClaims := token.Claims.(jwt.MapClaims)
	if !okClaims {
		return 0, fmt.Errorf("failed to extract claims from token")
	}

	idData, okId := claims[idStr]
	if !okId {
		return 0, fmt.Errorf("failed to extract id from claims")
	}

	id := 0
	switch idData.(type) {
	case int:
		id = idData.(int)
	case float64:
		var floatId = idData.(float64)
		id = common.Round(floatId)
	default:
		return 0, fmt.Errorf("failed to cast claims[id] to int")
	}

	return id, nil
}
