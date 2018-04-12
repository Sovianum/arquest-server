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
		c.JSON(http.StatusBadRequest, common.GetErrorJson(err))
		return
	}
	if code, err := validateUser(&user); err != nil {
		c.JSON(code, common.GetErrorJson(err))
	}

	exists, existsErr := env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(existsErr))
		return
	}
	if exists {
		c.JSON(http.StatusConflict, common.GetErrorJson(fmt.Errorf("user already exists")))
		return
	}

	hash, err := env.hashFunc([]byte(user.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(err))
		return
	}
	user.Password = string(hash)

	userId, saveErr := env.userDAO.Save(user)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(saveErr))
		return
	}

	tokenString, tokenErr := env.generateTokenString(userId, user.Login)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(tokenErr))
		// TODO add info that u has been successfully saved
		return
	}
	c.JSON(http.StatusOK, common.GetDataJson(tokenString))
}

func (env *Env) UserSignInPost(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.GetErrorJson(err))
		return
	}
	if code, err := validateUser(&user); err != nil {
		c.JSON(code, common.GetErrorJson(err))
	}

	exists, existsErr := env.userDAO.ExistsByLogin(user.Login)
	if existsErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(existsErr))
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, common.GetErrorJson(notFoundErr))
		return
	}

	dbUser, dbErr := env.userDAO.GetUserByLogin(user.Login)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(dbErr))
		return
	}

	if err := env.hashValidator([]byte(user.Password), []byte(dbUser.Password)); err != nil {
		c.JSON(http.StatusNotFound, common.GetErrorJson(err))
		return
	}

	tokenString, tokenErr := env.generateTokenString(dbUser.Id, dbUser.Login)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrorJson(tokenErr))
		return
		// TODO add info that u has been successfully saved
	}
	c.JSON(http.StatusOK, common.GetDataJson(tokenString))
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
	if user.Login == "" || user.Password == "" {
		return http.StatusBadRequest, fmt.Errorf("empty user")
	}
	return http.StatusOK, nil
}
