package httputils

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrGetUserId = errors.New("User Id not exists")
)

func GetUserId(ctx *gin.Context) (uint, error) {
	anyId, ok := ctx.Get("user_id")

	if !ok {
		return 0, ErrGetUserId
	}

	id, ok := anyId.(float64)

	if !ok {
		return 0, ErrGetUserId
	}

	return uint(id), nil
}
