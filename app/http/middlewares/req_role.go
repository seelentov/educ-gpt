package middlewares

import (
	"educ-gpt/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func RequiredRolesMiddleware(
	roleNames []string,
	logger *zap.Logger,
	roleService services.RoleService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		anyUserId, ok := c.Get("user_id")

		if !ok {
			logger.Error("user_id not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		userId, ok := anyUserId.(uint)

		if !ok {
			logger.Error("Cant parse uint from user id", zap.Any("user_id", anyUserId))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return // Important to return after aborting
		}

		roles, err := roleService.GetRolesByUserId(userId)
		if err != nil {
			logger.Error("Failed to retrieve roles for user", zap.Uint("user_id", userId), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return // Important to return after aborting
		}

		for _, requiredRole := range roleNames {
			for _, userRole := range roles {
				if userRole.Name == requiredRole {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}
