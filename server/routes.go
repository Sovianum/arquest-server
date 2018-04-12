package server

import (
	"github.com/gin-gonic/gin"
)

func GetEngine(env *Env) *gin.Engine {
	router := gin.Default()

	router.POST("/api/v1/auth/register", env.UserRegisterPost)
	router.POST("/api/v1/auth/login", env.UserSignInPost)
	router.GET("/api/v1/user/self", env.UserGetSelfInfo)

	return router
}
