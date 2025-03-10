package dtos

import "mime/multipart"

type UpdateUserRequest struct {
	Name         string                `form:"name,omitempty" json:"name,omitempty" binding:"omitempty,lte=100"`
	Number       string                `form:"number,omitempty" json:"number,omitempty" binding:"omitempty,gte=8,lte=100,number"`
	ChatGptModel string                `form:"chat_gpt_model,omitempty" json:"chat_gpt_model,omitempty" binding:"omitempty,lte=200"`
	ChatGptToken string                `form:"chat_gpt_token,omitempty" json:"chat_gpt_token,omitempty" binding:"omitempty,lte=200"`
	AvatarFile   *multipart.FileHeader `form:"avatar_file,omitempty" json:"avatar_file,omitempty" binding:"omitempty"`
}
