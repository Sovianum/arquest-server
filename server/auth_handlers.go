package server

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/model"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	idStr    = "id"
	loginStr = "login"
	expStr   = "exp"
)

var notFoundErr = fmt.Errorf("not found")

func (env *Env) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	env.logger.LogRequestStart(r)
	u, code, parseErr := parseUser(r)
	if parseErr != nil {
		env.logger.LogRequestError(r, parseErr)
		w.WriteHeader(code)
		common.WriteWithLogging(r, w, common.GetErrorJson(parseErr), env.logger)
		return
	}

	exists, existsErr := env.userDAO.ExistsByLogin(u.Login)
	if existsErr != nil {
		env.logger.LogRequestError(r, existsErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(existsErr), env.logger)
		return
	}
	if exists {
		err := fmt.Errorf("user already exists")
		env.logger.LogRequestError(r, err)
		w.WriteHeader(http.StatusConflict)
		common.WriteWithLogging(r, w, common.GetErrorJson(err), env.logger)
		return
	}

	hash, err := env.hashFunc([]byte(u.Password))
	if err != nil {
		env.logger.LogRequestError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(err), env.logger)
		return
	}
	u.Password = string(hash)

	userId, saveErr := env.userDAO.Save(u)
	if saveErr != nil {
		env.logger.LogRequestError(r, saveErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(saveErr), env.logger)
		return
	}

	tokenString, tokenErr := env.generateTokenString(userId, u.Login)
	if tokenErr != nil {
		env.logger.LogRequestError(r, tokenErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(tokenErr), env.logger)
		// TODO add info that u has been successfully saved
		return
	}

	env.logger.LogRequestSuccess(r)
	common.WriteWithLogging(r, w, common.GetDataJson(tokenString), env.logger)
}

func (env *Env) UserSignInPost(w http.ResponseWriter, r *http.Request) {
	env.logger.LogRequestStart(r)
	u, code, parseErr := parseUser(r)
	if parseErr != nil {
		env.logger.LogRequestError(r, parseErr)
		w.WriteHeader(code)
		common.WriteWithLogging(r, w, common.GetErrorJson(parseErr), env.logger)
		return
	}

	exists, existsErr := env.userDAO.ExistsByLogin(u.Login)
	if existsErr != nil {
		env.logger.LogRequestError(r, existsErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(existsErr), env.logger)
		return
	}
	if !exists {
		env.logger.LogRequestError(r, notFoundErr)
		w.WriteHeader(http.StatusNotFound)
		common.WriteWithLogging(r, w, common.GetErrorJson(notFoundErr), env.logger)
		return
	}

	dbUser, dbErr := env.userDAO.GetUserByLogin(u.Login)
	if dbErr != nil {
		env.logger.LogRequestError(r, dbErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(dbErr), env.logger)
		return
	}

	if err := env.hashValidator([]byte(u.Password), []byte(dbUser.Password)); err != nil {
		env.logger.LogRequestError(r, err)
		w.WriteHeader(http.StatusNotFound)
		common.WriteWithLogging(r, w, common.GetErrorJson(err), env.logger)
		return
	}

	tokenString, tokenErr := env.generateTokenString(dbUser.Id, dbUser.Login)
	if tokenErr != nil {
		env.logger.LogRequestError(r, tokenErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(tokenErr), env.logger)
		// TODO add info that u has been successfully saved
		return
	}

	env.logger.LogRequestSuccess(r)
	common.WriteWithLogging(r, w, common.GetDataJson(tokenString), env.logger)
}

func (env *Env) UserGetSelfInfo(w http.ResponseWriter, r *http.Request) {
	env.logger.LogRequestStart(r)
	userId, idCode, idErr := env.getIdFromRequest(r)
	if idErr != nil {
		env.logger.LogRequestError(r, idErr)
		w.WriteHeader(idCode)
		common.WriteWithLogging(r, w, common.GetErrorJson(idErr), env.logger)
		return
	}

	exists, existsErr := env.userDAO.ExistsById(userId)
	if existsErr != nil {
		env.logger.LogRequestError(r, existsErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(existsErr), env.logger)
		return
	}
	if !exists {
		env.logger.LogRequestError(r, notFoundErr)
		w.WriteHeader(http.StatusNotFound)
		common.WriteWithLogging(r, w, common.GetErrorJson(notFoundErr), env.logger)
		return
	}

	var dbUser, dbErr = env.userDAO.GetUserById(userId)
	if dbErr != nil {
		env.logger.LogRequestError(r, dbErr)
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteWithLogging(r, w, common.GetErrorJson(dbErr), env.logger)
		return
	}

	env.logger.LogRequestSuccess(r)
	common.WriteWithLogging(r, w, common.GetDataJson(dbUser), env.logger)
}

func (env *Env) generateTokenString(id int, login string) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	claims := t.Claims.(jwt.MapClaims)

	claims[idStr] = id
	claims[loginStr] = login
	claims[expStr] = time.Now().Add(time.Hour * 24 * time.Duration(env.conf.Auth.ExpireDays)).Unix()

	var tokenKey = env.conf.Auth.GetTokenKey()
	return t.SignedString(tokenKey)
}

func parseUser(r *http.Request) (*model.User, int, error) {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var u = new(model.User)
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if u.Login == "" || u.Password == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("empty user")
	}

	return u, http.StatusOK, nil
}
