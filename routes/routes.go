package routes

import (
	"github.com/Sovianum/arquest-server/server"
	"github.com/gin-gonic/gin"
)

func GetEngine(env *server.Env) *gin.Engine {
	router := gin.Default()

	root := router.Group("/api/v1/")
	root.GET("quests", env.GetAllQuests)

	authGroup := root.Group("auth")
	authGroup.POST("register", env.UserRegisterPost)
	authGroup.POST("login", env.UserSignInPost)

	userGroup := root.Group("user")
	userGroup.Use(env.CheckAuthorization)
	userGroup.GET("/", env.UserGetSelfInfo)

	questGroup := userGroup.Group("quest")
	questGroup.GET("finished", env.GetFinishedQuests)

	voteGroup := userGroup.Group("vote")
	voteGroup.GET("all", env.GetUserVotes)
	voteGroup.POST("mark", env.MarkQuest)
	voteGroup.POST("finish", env.FinishQuest)

	return router
}
