package middlewares

import (
	"educ-gpt/services"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func AuthMiddleware(logger *zap.Logger, jwtService services.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := jwtService.ValidateToken(tokenString)

		if err != nil {
			logger.Error("Invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID := claims["user_id"]
		c.Set("user_id", userID)
		c.Next()
	}
}
