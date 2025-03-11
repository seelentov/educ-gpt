package dtos

type ResetPasswordTaskRequest struct {
	Credential string `json:"credential" binding:"required"`
}
