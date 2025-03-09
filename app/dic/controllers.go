package dic

import (
	"educ-gpt/http/controllers"
	"educ-gpt/logger"
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
			ResetTokenService(),
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
