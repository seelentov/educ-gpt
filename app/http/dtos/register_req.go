package dtos

type RegisterRequest struct {
	Name         string `json:"name" binding:"required,lte=100"`
	Email        string `json:"email" binding:"required,email,lte=100"`
	Password     string `json:"password" binding:"required,gte=8"`
	ChatGptToken string `json:"chat_gpt_token" binding:"required"`
}
