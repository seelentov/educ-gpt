package services

import "errors"

var (
	ErrRequestFailed = errors.New("request failed")
	ErrParseFailed   = errors.New("request failed")
)

type NaturalLanguageService interface {
	GetAnswer(token string, model string, dialog []*DialogItem, target interface{}) error
}

type DialogItem struct {
	Text   string
	IsUser bool
}
