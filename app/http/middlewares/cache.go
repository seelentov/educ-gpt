package middlewares

import (
	"bytes"
	"educ-gpt/utils/httputils"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheMiddleware(logger *zap.Logger, cache *redis.Client, d time.Duration, onSession bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if onSession {
			userIdValue, exist := c.Get("user_id")

			if exist {
				userId, ok := userIdValue.(uint)
				if !ok {
					logger.Error("Cant parse user id")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
					return
				}

				path = fmt.Sprintf("%v!id=%v", path, userId)
			}
		}

		cachedData, err := cache.Get(c, path).Result()
		if err == nil {
			var data interface{}
			if err := json.Unmarshal([]byte(cachedData), &data); err == nil {
				c.JSON(200, data)
				c.Abort()
				return
			}
		}

		w := httputils.NewResponseBodyWriter(c.Writer, &bytes.Buffer{})
		c.Writer = w

		c.Next()

		if c.Writer.Status() == 200 {
			_, err := cache.Set(c, path, w.Body().String(), d).Result()
			if err != nil {
				logger.Error("Cant set cache ", zap.Error(err))
				return
			}
		}
	}
}
