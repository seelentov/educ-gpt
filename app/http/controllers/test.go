package controllers

import (
	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func (c TestController) Test(ctx *gin.Context) {
	ctx.String(200, "")
}

func NewTestController() *TestController {
	return &TestController{}
}
