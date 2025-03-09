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

type AuthController struct {
	userService   services.UserService
	jwtService    services.JwtService
	roleService   services.RoleService
	senderSrv     services.SenderService
	mailSrv       services.MailService
	resetTokenSrv services.ResetTokenService
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	token := req.ChatGptModel

	if token == "" {
		token = "gpt-4o-mini"
	}

	user := &models.User{
		ID:           0,
		Name:         req.Name,
		Email:        req.Email,
		Number:       req.Number,
		Password:     req.Password,
		ChatGptModel: req.ChatGptModel,
		ChatGptToken: req.ChatGptToken,
	}

	key, err := c.userService.Create(user)

	if err != nil {
		if errors.Is(err, services.ErrAlreadyExist) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
			return
		}

		if errors.Is(err, services.ErrDuplicate) {
			if errors.Is(err, services.ErrDuplicateEmail) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same email already exists"})
				return
			}

			if errors.Is(err, services.ErrDuplicateNumber) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same number already exists"})
				return
			}

			if errors.Is(err, services.ErrDuplicateName) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same name already exists"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	mail, err := c.mailSrv.ActivateMail(user.Name, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := c.senderSrv.SendMessage(user.Email, mail.Subject, mail.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "Письмо для активации аккаунта отправлено на " + user.Email})
}

func (c *AuthController) Activate(ctx *gin.Context) {
	key := ctx.Param("key")

	if err := c.userService.Activate(key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "Аккаунт успешно активирован"})
}

func (c *AuthController) Me(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := c.userService.GetByCredential(req.Credential)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong login or password"})
		return
	}

	if user.ActivateAt == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified"})
		return
	}

	err = c.userService.VerifyPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong login or password"})
		return
	}

	token, err := c.jwtService.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var req dtos.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken := req.RefreshToken

	claims, err := c.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	newToken, err := c.jwtService.GenerateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	//newRefreshToken, err := c.jwtService.GenerateRefreshToken(userID)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"token": newToken,
		//"refresh_token": newRefreshToken,
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req dtos.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if req.Password == req.OldPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Old password cannot be the same as password"})
		return
	}

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = c.userService.VerifyPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	if err := c.userService.ChangePassword(userId, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	key := ctx.Param("key")

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.resetTokenSrv.Verify(userId, key); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong reset token"})
		return
	}

	if err := c.userService.ChangePassword(userId, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (c *AuthController) ResetPasswordTask(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	key, err := c.resetTokenSrv.Create(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	mail, err := c.mailSrv.ResetMail(user.Name, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := c.senderSrv.SendMessage(user.Email, mail.Subject, mail.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "Письмо для восстановления аккаунта отправлено на " + user.Email})
}

func NewAuthController(
	userService services.UserService,
	jwtService services.JwtService,
	roleService services.RoleService,
	senderSrv services.SenderService,
	mailSrv services.MailService,
	resetTokenSrv services.ResetTokenService,
) *AuthController {
	return &AuthController{userService, jwtService, roleService, senderSrv, mailSrv, resetTokenSrv}
}
