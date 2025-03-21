package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/httputils"
	"educ-gpt/utils/httputils/valid"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type DialogController struct {
	dialogSrv services.DialogService
	userSrv   services.UserService
	aiSrv     services.AIService
}

// GetDialogs returns the current user's dialogs
// @Summary      Get current user's dialogs
// @Description  Returns the current user's dialogs based on the JWT token
// @Tags         dialog
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer <JWT token>"
// @Success      200 {array} models.Dialog "User`s dialogs"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /dialogs [get]
func (d DialogController) GetDialogs(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	dialogs, err := d.dialogSrv.GetDialogsByUserID(userid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dialogs)
}

// GetDialog returns dialog
// @Summary      Get dialog by id
// @Description  Returns the dialog after verification based on the JWT token
// @Tags         dialog
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer <JWT token>"
// @Success      200 {array} models.Dialog "Dialog"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      403 {object} dtos.ErrorResponse "Forbidden"
// @Failure      404 {object} dtos.ErrorResponse "Not found"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /dialogs/{dialog_id} [get]
func (d DialogController) GetDialog(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	dialogId, err := strconv.ParseUint(ctx.Param("dialog_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	dialog, err := d.dialogSrv.GetDialog(uint(dialogId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if dialog.UserID != userid {
		ctx.JSON(http.StatusForbidden, dtos.ForbiddenResponse())
		return
	}

	ctx.JSON(http.StatusOK, dialog)
}

// ThrowMessage to dialog and get answer
// @Summary      Add message to dialog and get answer
// @Description  Add message to dialog by id after verification based on the JWT token and get answer
// @Tags         dialog
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.ThrowMessageRequest true "Message"
// @Success      200 {array} models.DialogItem "AI response"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      403 {object} dtos.ErrorResponse "Forbidden"
// @Failure      404 {object} dtos.ErrorResponse "Not found"
// @Failure      409 {object} dtos.ErrorResponse "AI Error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /dialogs/{dialog_id} [post]
func (d DialogController) ThrowMessage(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	dialogId, err := strconv.ParseUint(ctx.Param("dialog_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	dialog, err := d.dialogSrv.GetDialog(uint(dialogId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if dialog.UserID != userid {
		ctx.JSON(http.StatusForbidden, dtos.ForbiddenResponse())
		return
	}

	var req dtos.ThrowMessageRequest

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

	err = d.dialogSrv.AddDialogItem(&models.DialogItem{
		Text:     req.Message,
		IsUser:   true,
		DialogID: uint(dialogId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	user, err := d.userSrv.GetById(userid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	dialogItems := dialog.DialogItems[:]

	dialogItems = append(dialogItems, &models.DialogItem{Text: req.Message, IsUser: true})

	var target string

	err = d.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, dialogItems, &target)
	if err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, services.ErrParseResFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Неверный формат ответа AI. Попробуйте еще раз."})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	message := &models.DialogItem{
		Text:      target,
		IsUser:    false,
		DialogID:  uint(dialogId),
		CreatedAt: time.Now(),
	}

	err = d.dialogSrv.AddDialogItem(message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// RemoveDialog
// @Description  Remove dialog by id after verification based on the JWT token and get answer
// @Tags         dialog
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer <JWT token>"
// @Success      200 {object} dtos.StatusResponse "Dialog removed successfully"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      403 {object} dtos.ErrorResponse "Forbidden"
// @Failure      404 {object} dtos.ErrorResponse "Not found"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /dialogs/{dialog_id} [delete]
func (d DialogController) RemoveDialog(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	dialogId, err := strconv.ParseUint(ctx.Param("dialog_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	dialog, err := d.dialogSrv.GetDialog(uint(dialogId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if dialog.UserID != userid {
		ctx.JSON(http.StatusForbidden, dtos.ForbiddenResponse())
		return
	}

	if err := d.dialogSrv.RemoveDialog(uint(dialogId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.OkResponse())
}

// CreateDialog
// @Summary      Create dialog
// @Description  Create dialog with user_id based on the JWT token
// @Tags         dialog
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer <JWT token>"
// @Success      200 {object} models.Dialog "Dialog created successfully"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /dialogs [post]
func (d DialogController) CreateDialog(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	dialog, err := d.dialogSrv.CreateDialog(&models.Dialog{
		UserID:      userid,
		DialogItems: []*models.DialogItem{{Text: "Привет. Что ты хочешь узнать?"}},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dialog)
}

func NewDialogController(
	dialogSrv services.DialogService,
	userSrv services.UserService,
	aiSrv services.AIService,
) *DialogController {
	return &DialogController{dialogSrv, userSrv, aiSrv}
}
