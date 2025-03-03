package httputils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrGetUserId = errors.New("User Id not exists")
)

func GetUserId(ctx *gin.Context) (uint, error) {
	anyId, ok := ctx.Get("user_id")

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0, ErrGetUserId
	}

	id, ok := anyId.(float64)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return 0, ErrGetUserId
	}

	return uint(id), nil
}
