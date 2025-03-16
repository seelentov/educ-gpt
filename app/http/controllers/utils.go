package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/httputils"
	"educ-gpt/utils/httputils/valid"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
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
// @Param Authorization header string true "Bearer <JWT token>"
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

	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := u.userSrv.GetById(userid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := u.promptSrv.CompileCode(req.Code, req.Language)

	var res dtos.ResultResponse

	err = u.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*models.DialogItem{{Text: prompt, IsUser: true}}, &res)
	if err != nil {
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
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.VerifyAnswerRequest true "Answer details"
// @Success      200 {object} services.PromptProblemResponse "Answer verification result"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /utils/check_answer [post]
func (r RoadmapController) VerifyAnswer(ctx *gin.Context) {
	var req dtos.VerifyAnswerRequest

	userId, err := httputils.GetUserId(ctx)
	if err != nil {

		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

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

	err = r.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*models.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
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
