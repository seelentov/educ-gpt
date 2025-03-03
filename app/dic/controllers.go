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
		)
		logger.Logger().Debug("AuthController initialized")
	}

	return authController
}

var testController *controllers.TestController

func TestController() *controllers.TestController {
	if testController == nil {
		testController = controllers.NewTestController()
		logger.Logger().Debug("TestController initialized")
	}

	return testController
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
