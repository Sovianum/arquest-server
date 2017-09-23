package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Sovianum/acquaintanceServer/model"
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

func (env *Env) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	var user, code, parseErr = parseUser(r)
	if parseErr != nil {
		w.WriteHeader(code)
		return
	}

	var exists, existsErr = env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	var userId, saveErr = env.userDAO.Save(user)
	if saveErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tokenString, tokenErr = env.getTokenString(userId, user.Login)
	if tokenErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO add info that user has been successfully saved
		return
	}

	w.Write([]byte(tokenString))
}

func (env *Env) UserSignInPost(w http.ResponseWriter, r *http.Request) {
	var user, code, parseErr = parseUser(r)
	if parseErr != nil {
		w.WriteHeader(code)
		return
	}

	var exists, existsErr = env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var userId, idErr = env.userDAO.GetIdByLogin(user.Login)
	if idErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tokenString, tokenErr = env.getTokenString(userId, user.Login)
	if tokenErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO add info that user has been successfully saved
		return
	}

	w.Write([]byte(tokenString))
}

func (env *Env) getTokenString(id int, login string) (string, error) {
	var token = jwt.New(jwt.SigningMethodHS256)
	var claims = token.Claims.(jwt.MapClaims)

	claims[idStr] = id
	claims[loginStr] = login
	claims[expStr] = time.Now().Add(time.Hour * 24 * time.Duration(env.authConf.ExpireDays)).Unix()

	var tokenKey = env.authConf.GetTokenKey()
	return token.SignedString(tokenKey)
}

func parseUser(r *http.Request) (*model.User, int, error) {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var user = new(model.User)
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if user.Login == "" || user.Password == "" {
		return nil, http.StatusBadRequest, errors.New("Empty user")
	}

	return user, http.StatusOK, nil
}
