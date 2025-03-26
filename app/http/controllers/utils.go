package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/httputils/valid"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UtilsController struct {
	aiSrv     services.AIService
	promptSrv services.PromptService
	userSrv   services.UserService
}

// Compile code by AI
// @Summary      Compile
// @Description  Compile code by AI
// @Tags         utils
// @Accept       json
// @Produce      json
// @Param        request body dtos.CompileRequest true "Code for compiler"
// @Success      200 {object} dtos.ResultResponse "Compiled code"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /utils/compile [post]
func (u UtilsController) Compile(ctx *gin.Context) {
	var req dtos.CompileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {

			ctx.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := u.promptSrv.CompileCode(req.Code, req.Language)

	var res dtos.ResultResponse

	if err := u.aiSrv.GetAnswer("", "", []*models.DialogItem{{Text: prompt, IsUser: true}}, &res); err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// VerifyAnswer and get verification status by AI
// @Summary      Verify answer
// @Description  VerifyAnswer and get verification status by AI
// @Tags         utils
// @Accept       json
// @Produce      json
// @Param        request body dtos.VerifyAnswerRequest true "Answer details"
// @Success      200 {object} services.PromptProblemResponse "Answer verification result"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /utils/check_answer [post]
func (r UtilsController) VerifyAnswer(ctx *gin.Context) {
	var req dtos.VerifyAnswerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := r.promptSrv.VerifyAnswer(req.Problem, req.Answer, req.Language)
	var target services.PromptProblemResponse

	if err := r.aiSrv.GetAnswer("", "", []*models.DialogItem{{Text: prompt, IsUser: true}}, &target); err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, target)
}

func NewUtilsController(
	aiSrv services.AIService,
	promptSrv services.PromptService,
	userSrv services.UserService,
) *UtilsController {
	return &UtilsController{aiSrv, promptSrv, userSrv}
}
