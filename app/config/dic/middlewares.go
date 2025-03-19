package dic

import (
	"educ-gpt/config/data"
	"educ-gpt/config/logger"
	"educ-gpt/http/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

var authMiddleware gin.HandlerFunc

func AuthMiddleware() gin.HandlerFunc {
	if authMiddleware == nil {
		authMiddleware = middlewares.AuthMiddleware(
			logger.Logger(),
			os.Getenv("JWT_SECRET"),
		)
		log.Print("AuthMiddleware initialized")
	}

	return authMiddleware
}

var requiredAuthMiddleware gin.HandlerFunc

func RequiredAuthMiddleware() gin.HandlerFunc {
	if requiredAuthMiddleware == nil {
		requiredAuthMiddleware = middlewares.RequiredAuthMiddleware()
		log.Print("RequiredAuthMiddleware initialized")
	}

	return requiredAuthMiddleware
}

func RequiredRoleMiddleware(roleNames []string) gin.HandlerFunc {
	return middlewares.RequiredRolesMiddleware(
		roleNames,
		logger.Logger(),
		RoleService(),
	)
}

func CacheMiddleware(d time.Duration, onSession bool) gin.HandlerFunc {
	cacheDefaultDuration, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP"))

	if err != nil {
		log.Fatal(err)
	}

	if d == 0 {
		d = time.Duration(cacheDefaultDuration) * time.Second
	}

	return middlewares.CacheMiddleware(
		logger.Logger(),
		data.Redis(),
		d,
		onSession,
	)
}
