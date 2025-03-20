package dic

import (
	"educ-gpt/config/data"
	"educ-gpt/config/logger"
	"educ-gpt/services"
	"educ-gpt/services/impl"
	"log"
	"os"
	"strconv"
)

var userService services.UserService

func UserService() services.UserService {
	if userService == nil {
		userService = impl.NewUserServiceImpl(
			data.DB(),
			logger.Logger(),
			os.Getenv("AUTH_DEFAULT_ROLE"),
		)
		log.Print("UserService initialized")
	}

	return userService
}

var roleService services.RoleService

func RoleService() services.RoleService {
	if roleService == nil {
		roleService = impl.NewRoleServiceImpl(
			data.DB(),
			logger.Logger(),
			os.Getenv("AUTH_DEFAULT_ROLE"),
		)
		log.Print("RoleService initialized")
	}

	return roleService
}

var roadmapService services.RoadmapService

func RoadmapService() services.RoadmapService {
	if roadmapService == nil {
		roadmapService = impl.NewRoadmapServiceImpl(
			data.DB(),
			logger.Logger(),
		)
		log.Print("RoadmapService initialized")
	}

	return roadmapService
}

var promptService services.PromptService

func PromptService() services.PromptService {
	if promptService == nil {
		promptService = impl.NewPromptServiceImpl()
		log.Print("PromptService initialized")
	}

	return promptService
}

var jwtService services.JwtService

func JwtService() services.JwtService {
	if jwtService == nil {

		jwtExpiration, err := strconv.Atoi(os.Getenv("JWT_EXP"))

		if err != nil {
			log.Fatal(err)
		}

		jwtRefreshExpiration, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP"))

		if err != nil {
			log.Fatal(err)
		}

		jwtService = impl.NewJwtServiceImpl(
			os.Getenv("JWT_SECRET"),
			os.Getenv("JWT_REFRESH_SECRET"),
			jwtExpiration,
			jwtRefreshExpiration,
			logger.Logger(),
		)
		log.Print("JwtService initialized")
	}

	return jwtService
}

var gptService services.GptService

func GptService() services.GptService {
	if gptService == nil {
		gptService = impl.NewGptService(
			logger.Logger(),
		)
		log.Print("GptService initialized")
	}

	return gptService
}

var senderService services.SenderService

func SenderService() services.SenderService {
	smtpPort, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	if senderService == nil {
		senderService = impl.NewSenderServiceImpl(
			os.Getenv("SMTP_HOST"),
			int(smtpPort),
			os.Getenv("SMTP_USERNAME"),
			os.Getenv("SMTP_PASSWORD"),
			os.Getenv("SMTP_FROM"),
			"email_queue",
			logger.Logger(),
			data.Redis(),
		)
		log.Print("SenderService initialized")
	}

	return senderService
}

var mailService services.MailService

func MailService() services.MailService {
	if mailService == nil {
		mailService = impl.NewMailServiceImpl(
			os.Getenv("PROTOCOL"),
			os.Getenv("FULL_HOST"),
			"activate",
			"reset",
			"change_email",
			logger.Logger(),
		)
		log.Print("MailService initialized")
	}

	return mailService
}

var tokenService services.TokenService

func TokenService() services.TokenService {
	if tokenService == nil {
		tokenService = impl.NewTokenServiceImpl(
			data.DB(),
			logger.Logger(),
		)
		log.Print("TokenService initialized")
	}

	return tokenService
}

var fileService services.FileService

func FileService() services.FileService {
	if fileService == nil {
		fileService = impl.NewFileServiceImpl(
			logger.Logger(),
			"/storage",
		)
		log.Print("FileService initialized")
	}

	return fileService
}

var dialogService services.DialogService

func DialogService() services.DialogService {
	if dialogService == nil {
		dialogService = impl.NewDialogService(
			data.DB(),
			logger.Logger(),
		)
		log.Print("DialogService initialized")
	}

	return dialogService
}
