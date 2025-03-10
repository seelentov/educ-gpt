package dtos

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required,email,lte=100"`
}
