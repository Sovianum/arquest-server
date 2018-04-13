package server

import (
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	idStr    = "id"
	loginStr = "login"
	expStr   = "exp"
)

var notFoundErr = fmt.Errorf("not found")

func (env *Env) UserRegisterPost(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.GetErrResponse(err))
		return
	}
	if code, err := validateUser(&user); err != nil {
		c.JSON(code, common.GetErrResponse(err))
		return
	}

	exists, existsErr := env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(existsErr))
		return
	}
	if exists {
		c.JSON(http.StatusConflict, common.GetErrResponse(fmt.Errorf("user already exists")))
		return
	}

	hash, err := env.hashFunc([]byte(user.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	user.Password = string(hash)

	userId, saveErr := env.userDAO.Save(user)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(saveErr))
		return
	}

	tokenString, tokenErr := env.generateTokenString(userId, user.Login)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(tokenErr))
		// TODO add info that u has been successfully saved
		return
	}
	c.JSON(http.StatusOK, common.GetDataResponse(tokenString))
}

func (env *Env) UserSignInPost(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.GetErrResponse(err))
		return
	}
	if code, err := validateUser(&user); err != nil {
		c.JSON(code, common.GetErrResponse(err))
	}

	exists, existsErr := env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(existsErr))
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, common.GetErrResponse(notFoundErr))
		return
	}

	dbUser, dbErr := env.userDAO.GetUserByLogin(user.Login)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(dbErr))
		return
	}

	if err := env.hashValidator([]byte(user.Password), []byte(dbUser.Password)); err != nil {
		c.JSON(http.StatusNotFound, common.GetErrResponse(err))
		return
	}

	tokenString, tokenErr := env.generateTokenString(dbUser.Id, dbUser.Login)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(tokenErr))
		return
		// TODO add info that u has been successfully saved
	}
	c.JSON(http.StatusOK, common.GetDataResponse(tokenString))
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

func validateUser(user *model.User) (int, error) {
	if user.Login == "" {
		return http.StatusBadRequest, fmt.Errorf("no login")
	}
	if user.Password == "" {
		return http.StatusBadRequest, fmt.Errorf("no password")
	}
	return http.StatusOK, nil
}
