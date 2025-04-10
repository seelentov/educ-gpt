package dic

import (
	"educ-gpt/http/controllers"
	"log"
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
			DialogService(),
		)
		log.Print("AuthController initialized")
	}

	return authController
}

var roadmapController *controllers.RoadmapController

func RoadmapController() *controllers.RoadmapController {
	if roadmapController == nil {
		roadmapController = controllers.NewRoadmapController(
			UserService(),
			AIService(),
			PromptService(),
			RoadmapService(),
		)
		log.Print("RoadmapController initialized")
	}

	return roadmapController
}

var utilsController *controllers.UtilsController

func UtilsController() *controllers.UtilsController {
	if utilsController == nil {
		utilsController = controllers.NewUtilsController(
			AIService(),
			PromptService(),
			UserService(),
		)
		log.Print("UtilsController initialized")
	}
	return utilsController
}

var dialogController *controllers.DialogController

func DialogController() *controllers.DialogController {
	if dialogController == nil {
		dialogController = controllers.NewDialogController(
			DialogService(),
			UserService(),
			AIService(),
		)
		log.Print("DialogController initialized")
	}
	return dialogController
}
