package server

import (
	"github.com/Sovianum/arquest-server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserID = "userID"
)

func (env *Env) CheckAuthorization(c *gin.Context) {
	userId, idCode, idErr := env.getIdFromRequest(c.Request)
	if idErr != nil {
		c.JSON(idCode, common.GetErrResponse(idErr))
		c.Abort()
		return
	}

	exists, existsErr := env.userDAO.ExistsById(userId)
	if existsErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(existsErr))
		c.Abort()
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, common.GetErrResponse(notFoundErr))
		c.Abort()
		return
	}
	c.Set(UserID, userId)
	c.Next()
}
