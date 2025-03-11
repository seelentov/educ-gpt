package dic

import (
	"educ-gpt/config/logger"
	"educ-gpt/http/controllers"
)

var authController *controllers.AuthController

func AuthController() *controllers.AuthController {
	if authController == nil {
		authController = controllers.NewAuthController(
			UserService(),
			JwtService(),
			RoleService(),
			SenderService(),
			MailService(),
			TokenService(),
			FileService(),
		)
		logger.Logger().Debug("AuthController initialized")
	}

	return authController
}

var roadmapController *controllers.RoadmapController

func RoadmapController() *controllers.RoadmapController {
	if roadmapController == nil {
		roadmapController = controllers.NewRoadmapController(
			UserService(),
			GptService(),
			PromptService(),
			RoadmapService(),
		)
		logger.Logger().Debug("TestController initialized")
	}

	return roadmapController
}
