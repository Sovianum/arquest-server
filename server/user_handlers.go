package server

import (
	"github.com/Sovianum/arquest-server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) UserGetSelfInfo(c *gin.Context) {
	userId := c.GetInt(UserID)
	var dbUser, dbErr = env.userDAO.GetUserById(userId)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(dbErr))
		return
	}
	c.JSON(http.StatusOK, common.GetDataResponse(dbUser))
}
