package services

import (
	"errors"
)

var (
	ErrCreateResetToken = errors.New("cannot create a reset token")
	ErrVerifyResetToken = errors.New("cannot verify the reset token")
	ErrClearResetTokens = errors.New("cannot clear the reset tokens")
)

type ResetTokenService interface {
	Create(userID uint) (string, error)
	Verify(userID uint, key string) error
	Clear() error
}
