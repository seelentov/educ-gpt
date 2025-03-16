package router

import (
	"educ-gpt/config/dic"
	_ "educ-gpt/docs"
	"educ-gpt/http/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(dic.AuthMiddleware())

	router.Static("/storage", "./storage")

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
				authGroup.POST("/register", dic.AuthController().Register)
				authGroup.POST("/me", dic.RequiredAuthMiddleware(), dic.AuthController().Me)
				authGroup.POST("/login", dic.AuthController().Login)
				authGroup.POST("/refresh", dic.AuthController().Refresh)
				authGroup.POST("/activate/:key", dic.AuthController().Activate)
				authGroup.POST("/change_password", dic.AuthController().ChangePassword)
				authGroup.POST("/reset/:key/:user_id", dic.AuthController().ResetPassword)
				authGroup.POST("/reset/task", dic.AuthController().ResetPasswordTask)
				authGroup.POST("/change_email/task", dic.AuthController().ChangeEmailTask)
				authGroup.POST("/change_email/:key/:user_id", dic.AuthController().ChangeEmail)
				authGroup.PATCH("/update", dic.AuthController().UpdateUser)
			}

			roadmapGroup := v1.Group("/roadmap")
			{
				roadmapGroup.GET("", dic.RoadmapController().GetTopics)
				roadmapGroup.GET("/:topic_id", dic.RoadmapController().GetThemes)
				roadmapGroup.GET("/info/topic/:topic_id/", dic.RoadmapController().GetTopicInfo)
				roadmapGroup.GET("/info/theme/:theme_id/", dic.RoadmapController().GetThemeInfo)
				roadmapGroup.GET("/:topic_id/:theme_id", dic.RoadmapController().GetTheme)
				roadmapGroup.GET("/problems/:topic_id/:theme_id", dic.RoadmapController().GetProblems)
				roadmapGroup.POST("/resolve", dic.RoadmapController().VerifyAnswerAndIncrementUserScore)
			}

			utilsGroup := v1.Group("/utils")
			{
				utilsGroup.POST("/compile", dic.UtilsController().Compile)
				utilsGroup.POST("/check_answer", dic.RoadmapController().VerifyAnswer)
			}

			dialogGroup := v1.Group("/dialogs")
			{
				dialogGroup.GET("/", dic.DialogController().GetDialogs)
				dialogGroup.POST("/", dic.DialogController().CreateDialog)
				dialogGroup.GET("/:dialog_id", dic.DialogController().GetDialog)
				dialogGroup.POST("/:dialog_id", dic.DialogController().ThrowMessage)
				dialogGroup.DELETE("/:dialog_id", dic.DialogController().RemoveDialog)
			}
		}
	}

	return router
}
