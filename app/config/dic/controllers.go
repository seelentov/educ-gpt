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
		logger.Logger().Debug("RoadmapController initialized")
	}

	return roadmapController
}

var utilsController *controllers.UtilsController

func UtilsController() *controllers.UtilsController {
	if utilsController == nil {
		utilsController = controllers.NewUtilsController(
			GptService(),
			PromptService(),
			UserService(),
		)
		logger.Logger().Debug("UtilsController initialized")
	}
	return utilsController
}
