package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/httputils"
	"educ-gpt/utils/httputils/valid"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	userService services.UserService
	jwtService  services.JwtService
	roleService services.RoleService
	senderSrv   services.SenderService
	mailSrv     services.MailService
	tokenSrv    services.TokenService
}

// Register a new user
// @Summary      Register a new user
// @Description  Registers a new user and sends an activation email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RegisterRequest true "User registration details"
// @Success      201 {object} dtos.MessageResponse "Successfully registered, activation email sent"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      409 {object} dtos.ErrorResponse "User already exists"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest

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

	user := &models.User{
		ID:           0,
		Name:         req.Name,
		Email:        req.Email,
		Number:       req.Number,
		Password:     req.Password,
		ChatGptModel: "gpt-4o-mini",
		ChatGptToken: req.ChatGptToken,
	}

	key, err := c.userService.Create(user)

	if err != nil {
		if errors.Is(err, services.ErrAlreadyExist) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Пользователь уже существует"})
			return
		}

		if errors.Is(err, services.ErrDuplicate) {
			if errors.Is(err, services.ErrDuplicateEmail) {
				ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Пользователь с таким email уже существует"})
				return
			}

			if errors.Is(err, services.ErrDuplicateNumber) {
				ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Пользователь с таким номером телефона уже существует"})
				return
			}

			if errors.Is(err, services.ErrDuplicateName) {
				ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Пользователь с таким именем телефона уже существует"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	mail, err := c.mailSrv.ActivateMail(user.Name, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := c.senderSrv.SendMessage(user.Email, mail.Subject, mail.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusCreated, dtos.MessageResponse{Message: "Письмо для активации аккаунта отправлено на " + user.Email})
}

// Activate activates a user account
// @Summary      Activate a user account
// @Description  Activates a user account using the activation key
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        key path string true "Activation key"
// @Success      200 {object} dtos.StatusResponse "Account successfully activated"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/activate/{key} [post]
func (c *AuthController) Activate(ctx *gin.Context) {
	key := ctx.Param("key")

	if err := c.userService.Activate(key); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.StatusResponse{Status: "Аккаунт успешно активирован"})
}

// Me returns the current user's information
// @Summary      Get current user's information
// @Description  Returns the current user's information based on the JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {object} models.User "User information"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/me [post]
func (c *AuthController) Me(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

// Login authenticates a user and returns JWT tokens
// @Summary      Authenticate a user
// @Description  Authenticates a user and returns JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.LoginRequest true "User credentials"
// @Success      200 {object} dtos.TokenResponse "JWT tokens"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Invalid credentials or account not activated"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest

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

	user, err := c.userService.GetByCredential(req.Credential)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный логин или пароль"})
		return
	}

	if user.ActivateAt == nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Аккаунт не активирован"})
		return
	}

	err = c.userService.VerifyPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный логин или пароль"})
		return
	}

	token, err := c.jwtService.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.TokenResponse{Token: token, RefreshToken: refreshToken})
}

// Refresh refreshes the JWT token using a refresh token
// @Summary      Refresh JWT token
// @Description  Refreshes the JWT token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} dtos.TokenResponse "New JWT token"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Invalid refresh token"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req dtos.RefreshTokenRequest
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

	refreshToken := req.RefreshToken

	claims, err := c.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный refresh_token"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	newToken, err := c.jwtService.GenerateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.TokenResponse{Token: newToken})
}

// UpdateUser updates user information
// @Summary      Update user information
// @Description  Updates user information based on the provided data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.UpdateUserRequest true "User information to update"
// @Success      200 {object} dtos.StatusResponse "User information updated successfully"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/update [patch]
func (c *AuthController) UpdateUser(ctx *gin.Context) {
	var req dtos.UpdateUserRequest
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

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	updates := make(map[string]interface{})

	reqJson, err := json.Marshal(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := json.Unmarshal(reqJson, &updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := c.userService.Update(userId, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.OkResponse())
}

// ChangePassword changes the user's password
// @Summary      Change user's password
// @Description  Changes the user's password after verifying the old password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.ChangePasswordRequest true "Password change details"
// @Success      200 {object} dtos.StatusResponse "Password changed successfully"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized or invalid old password"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/change_password [post]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req dtos.ChangePasswordRequest

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

	if req.Password == req.OldPassword {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: "Новый пароль не должен совпадать со старым"})
		return
	}

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	err = c.userService.VerifyPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный пароль"})
		return
	}

	if err := c.userService.ChangePassword(userId, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.OkResponse())
}

// ChangeEmail changes the user's email
// @Summary      Change user's email
// @Description  Changes the user's email after verifying the key
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        key path string true "Change email key"
// @Success      200 {object} dtos.StatusResponse "Email changed successfully"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized or invalid key"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/change_email/{key} [post]
func (c *AuthController) ChangeEmail(ctx *gin.Context) {
	key := ctx.Param("key")

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	email, err := c.tokenSrv.VerifyAndGetData(userId, key, models.TypeChangeEmail)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный ключ"})
		return
	}

	updates := make(map[string]interface{})
	updates["email"] = email

	if err := c.userService.Update(userId, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}
	ctx.JSON(http.StatusOK, dtos.OkResponse())
}

// ChangeEmailTask initiates the process of changing the user's email
// @Summary      Initiate email change
// @Description  Initiates the process of changing the user's email by sending a verification email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.ChangeEmailRequest true "New email address"
// @Success      200 {object} dtos.MessageResponse "Verification email sent"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      409 {object} dtos.ErrorResponse "New email is the same as the old one or already in use"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/change_email/task [post]
func (c *AuthController) ChangeEmailTask(ctx *gin.Context) {
	var req dtos.ChangeEmailRequest
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

	userByEmail, err := c.userService.GetByEmail(req.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if userByEmail != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: "Пользователь с таким email уже существует"})
		return
	}

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	if user.Email == req.Email {
		ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: "Новый email не должен быть идентичен старому"})
		return
	}
	key, err := c.tokenSrv.Create(userId, models.TypeChangeEmail, req.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	mail, err := c.mailSrv.ChangeEmailMail(user.Name, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := c.senderSrv.SendMessage(req.Email, mail.Subject, mail.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.MessageResponse{Message: "Письмо для смены почты отправлено на " + req.Email})
}

// ResetPassword resets the user's password
// @Summary      Reset user's password
// @Description  Resets the user's password using a valid reset key
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        key path string true "Reset password key"
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.ResetPasswordRequest true "New password"
// @Success      200 {object} dtos.StatusResponse "Password reset successfully"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized or invalid key"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/reset/{key} [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dtos.ResetPasswordRequest
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

	key := ctx.Param("key")

	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	if err := c.tokenSrv.Verify(userId, key, models.TypeResetPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "Неверный ключ"})
		return
	}

	if err := c.userService.ChangePassword(userId, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.OkResponse())
}

// ResetPasswordTask initiates the process of resetting the user's password
// @Summary      Initiate password reset
// @Description  Initiates the process of resetting the user's password by sending a reset email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {object} dtos.MessageResponse "Reset email sent"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /auth/reset/task [post]
func (c *AuthController) ResetPasswordTask(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := c.userService.GetById(userId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	key, err := c.tokenSrv.Create(userId, models.TypeResetPassword, "")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	mail, err := c.mailSrv.ResetMail(user.Name, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := c.senderSrv.SendMessage(user.Email, mail.Subject, mail.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.MessageResponse{Message: "Письмо для восстановления аккаунта отправлено на " + user.Email})
}

func NewAuthController(
	userService services.UserService,
	jwtService services.JwtService,
	roleService services.RoleService,
	senderSrv services.SenderService,
	mailSrv services.MailService,
	tokenSrv services.TokenService,
) *AuthController {
	return &AuthController{userService, jwtService, roleService, senderSrv, mailSrv, tokenSrv}
}
