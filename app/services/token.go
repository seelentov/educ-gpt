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
	Create(userID uint, t models.Type, data string) (string, error)
	Verify(userID uint, key string, t models.Type) error
	VerifyAndGetData(userID uint, key string, t models.Type) (string, error)
	Clear() error
}
