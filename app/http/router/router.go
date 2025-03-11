package router

import (
	dic2 "educ-gpt/config/dic"
	_ "educ-gpt/docs"
	"educ-gpt/http/dtos"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(dic2.AuthMiddleware())

	router.Static("/storage", "./app/storage")

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
		return
	})

	apiGroup := router.Group("/api")
	{
		v1 := apiGroup.Group("/v1")
		{
			authGroup := v1.Group("/auth")
			{
				authGroup.POST("/register", dic2.AuthController().Register)
				authGroup.POST("/me", dic2.RequiredAuthMiddleware(), dic2.AuthController().Me)
				authGroup.POST("/login", dic2.AuthController().Login)
				authGroup.POST("/refresh", dic2.AuthController().Refresh)
				authGroup.POST("/activate/:key", dic2.AuthController().Activate)
				authGroup.POST("/change_password", dic2.AuthController().ChangePassword)
				authGroup.POST("/reset/:key/:user_id", dic2.AuthController().ResetPassword)
				authGroup.POST("/reset/task", dic2.AuthController().ResetPasswordTask)
				authGroup.POST("/change_email/task", dic2.AuthController().ChangeEmailTask)
				authGroup.POST("/change_email/:key/:user_id", dic2.AuthController().ChangeEmail)
				authGroup.PATCH("/update", dic2.AuthController().UpdateUser)
			}

			roadmapGroup := v1.Group("/roadmap")
			{
				roadmapGroup.GET("", dic2.RoadmapController().GetTopics)
				roadmapGroup.GET("/:topic_id", dic2.RoadmapController().GetThemes)
				roadmapGroup.GET("/:topic_id/info", dic2.RoadmapController().GetTopicInfo)
				roadmapGroup.GET("/:topic_id/:theme_id", dic2.RoadmapController().GetTheme)
				roadmapGroup.GET("/problems/:topic_id/:theme_id", dic2.RoadmapController().GetProblems)
				roadmapGroup.POST("/resolve", dic2.RoadmapController().IncrementUserScoreAndAddAnswer)
			}
		}
	}

	return router
}
