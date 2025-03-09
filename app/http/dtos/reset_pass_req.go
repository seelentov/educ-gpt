package dtos

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
