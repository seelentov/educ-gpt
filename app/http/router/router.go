package router

import (
	"educ-gpt/dic"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.Use(dic.AuthMiddleware())

	apiGroup := router.Group("/api")
	{
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/register", dic.AuthController().Register)
			authGroup.POST("/me", dic.RequiredAuthMiddleware(), dic.AuthController().Me)
			authGroup.POST("/login", dic.AuthController().Login)
			authGroup.POST("/refresh", dic.AuthController().Refresh)
		}

		testGroup := apiGroup.Group("/test")
		{
			testGroup.GET("", dic.TestController().Test)
		}

		roadmapGroup := apiGroup.Group("/roadmap")
		{
			roadmapGroup.GET("", dic.RoadmapController().GetTopics)
			roadmapGroup.GET("/:topic_id", dic.RoadmapController().GetThemes)
			roadmapGroup.GET("/:topic_id/:theme_id", dic.RoadmapController().GetTheme)
			roadmapGroup.POST("/resolve", dic.RoadmapController().IncrementUserScoreAndAddAnswer)
		}
	}

	return router
}
