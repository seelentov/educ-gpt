package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrCreateToken = errors.New("cannot create a token")
	ErrVerifyToken = errors.New("cannot verify the token")
	ErrClearTokens = errors.New("cannot clear the tokens")
)

type TokenService interface {
	Create(userID uint, t models.TokenType, data string) (string, error)
	Verify(userID uint, key string, t models.TokenType) error
	VerifyAndGetData(userID uint, key string, t models.TokenType) (string, error)
	Clear() error
}
