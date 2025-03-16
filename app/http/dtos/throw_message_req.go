package dtos

type ThrowMessageRequest struct {
	Message string `json:"message" binding:"required"`
}
