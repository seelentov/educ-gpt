package dic

import (
	"educ-gpt/data"
	"educ-gpt/logger"
	"educ-gpt/services"
	"log"
	"os"
	"strconv"
)

var userService services.UserService

func UserService() services.UserService {
	if userService == nil {
		userService = services.NewUserServiceImpl(
			data.DB(),
			logger.Logger(),
			os.Getenv("AUTH_DEFAULT_ROLE"),
		)
		logger.Logger().Debug("UserService initialized")
	}

	return userService
}

var roleService services.RoleService

func RoleService() services.RoleService {
	if roleService == nil {
		roleService = services.NewRoleServiceImpl(
			data.DB(),
			logger.Logger(),
			os.Getenv("AUTH_DEFAULT_ROLE"),
		)
		logger.Logger().Debug("RoleService initialized")
	}

	return roleService
}

var roadmapService services.RoadmapService

func RoadmapService() services.RoadmapService {
	if roadmapService == nil {
		roadmapService = services.NewRoadmapServiceImpl(
			data.DB(),
			logger.Logger(),
		)
		logger.Logger().Debug("RoadmapService initialized")
	}

	return roadmapService
}

var promptService services.PromptService

func PromptService() services.PromptService {
	if promptService == nil {
		promptService = services.NewPromptServiceImpl()
		logger.Logger().Debug("PromptService initialized")
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

		jwtService = services.NewJwtServiceImpl(
			os.Getenv("JWT_SECRET"),
			os.Getenv("JWT_REFRESH_SECRET"),
			jwtExpiration,
			jwtRefreshExpiration,
			logger.Logger(),
		)
		logger.Logger().Debug("JwtService initialized")
	}

	return jwtService
}

var gptService services.GptService

func GptService() services.GptService {
	if gptService == nil {
		gptService = services.NewGptService(
			logger.Logger(),
		)
		logger.Logger().Debug("GptService initialized")
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
		senderService = services.NewSenderServiceImpl(
			os.Getenv("SMTP_HOST"),
			int(smtpPort),
			os.Getenv("SMTP_USERNAME"),
			os.Getenv("SMTP_PASSWORD"),
			os.Getenv("SMTP_FROM"),
			logger.Logger(),
		)
		logger.Logger().Debug("SenderService initialized")
	}

	return senderService
}

var mailService services.MailService

func MailService() services.MailService {
	if mailService == nil {
		mailService = services.NewMailServiceImpl(
			os.Getenv("PROTOCOL"),
			os.Getenv("FULL_HOST"),
			"activate",
			"reset",
			"change_email",
			logger.Logger(),
		)
		logger.Logger().Debug("MailService initialized")
	}

	return mailService
}

var tokenService services.TokenService

func TokenService() services.TokenService {
	if tokenService == nil {
		tokenService = services.NewTokenServiceImpl(
			data.DB(),
			logger.Logger(),
		)
		logger.Logger().Debug("TokenService initialized")
	}

	return tokenService
}

var fileService services.FileService

func FileService() services.FileService {
	if fileService == nil {
		fileService = services.NewFileServiceImpl(
			logger.Logger(),
			"/storage",
		)
		logger.Logger().Debug("FileService initialized")
	}

	return fileService
}
