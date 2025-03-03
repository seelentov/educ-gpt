package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequiredAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("user_id")

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		c.Next()
	}
}
