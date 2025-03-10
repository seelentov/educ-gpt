package services

import "errors"

var (
	ErrRequestFailed   = errors.New("request failed")
	ErrParseFailed     = errors.New("request failed")
	ErrAIRequestFailed = errors.New("ai request failed")
)

type AIService interface {
	GetAnswer(token string, model string, dialog []*DialogItem, target interface{}) error
}

type DialogItem struct {
	Text   string
	IsUser bool
}
