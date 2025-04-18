package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrRequestFailed   = errors.New("request failed")
	ErrParseFailed     = errors.New("parse request failed")
	ErrParseResFailed  = errors.New("parse response failed")
	ErrAIRequestFailed = errors.New("ai request failed")
)

type AIService interface {
	GetAnswer(token string, model string, dialog []*models.DialogItem, target interface{}) error
}

type AIResponse struct {
	Choices []*AIResponseChoice `json:"choices"`
}

type AIResponseChoice struct {
	Message *AIResponseChoiceMessage `json:"message"`
}

type AIResponseChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
