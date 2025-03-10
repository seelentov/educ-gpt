package dtos

type UpdateUserRequest struct {
	Name         string `json:"name" binding:"lte=100"`
	Number       string `json:"number" binding:"gte=8,lte=100,number"`
	ChatGptModel string `json:"chat_gpt_model"`
	ChatGptToken string `json:"chat_gpt_token"`
	AvatarUrl    string `json:"avatar_url" binding:"url"`
}
