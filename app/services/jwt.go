package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrGenerateToken           = errors.New("error generating token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrParseToken              = errors.New("error parsing token")
	ErrInvalidToken            = errors.New("invalid token")
)

type JwtService interface {
	GenerateToken(userID uint, ttl int) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateRefreshToken(tokenString string) (jwt.MapClaims, error)
}
