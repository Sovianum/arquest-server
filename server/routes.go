package server

import (
	"github.com/gin-gonic/gin"
)

func GetEngine(env *Env) *gin.Engine {
	router := gin.Default()

	root := router.Group("/api/v1/")

	auth := root.Group("auth")
	auth.POST("register", env.UserRegisterPost)
	auth.POST("login", env.UserSignInPost)

	authorized := root.Group("user")
	authorized.Use(env.CheckAuthorization)

	authorized.GET("self", env.UserGetSelfInfo)
	return router
}
