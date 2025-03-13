package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrHashingPassword    = errors.New("error hashing password")
	ErrCreatingUser       = errors.New("error creating user")
	ErrRetrievingUser     = errors.New("error retrieving user")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordComparison = errors.New("error comparing password hashes")
	ErrAlreadyExist       = errors.New("user already exists")
	ErrUpdateUser         = errors.New("error updating user")
	ErrDuplicateName      = errors.New("duplicate name")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicate          = errors.New("duplicate")
	ErrActivate           = errors.New("error activate user")
	ErrDeleteUsers        = errors.New("error deleting user")
)

type UserService interface {
	Create(*models.User) (string, error)
	Update(id uint, data map[string]interface{}) error
	GetById(uint) (*models.User, error)
	GetByName(string) (*models.User, error)
	GetByEmail(string) (*models.User, error)
	GetByCredential(string) (*models.User, error)
	Verify(input string, credential string) error
	VerifyPassword(input string, password string) error
	Activate(key string) error
	ChangePassword(userId uint, password string) error
	ClearNonActivatedUsers() error
}
