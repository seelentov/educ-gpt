package impl

import (
	"educ-gpt/services"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type JwtServiceImpl struct {
	jwtSecret            string
	jwtRefreshSecret     string
	jwtExpiration        int
	jwtRefreshExpiration int

	logger *zap.Logger
}

func (j JwtServiceImpl) GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.jwtExpiration) * time.Second)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.jwtSecret))
	if err != nil {
		j.logger.Error("Error generating token", zap.Error(err))
		return "", fmt.Errorf("%w:%w", services.ErrGenerateToken, err)
	}

	return tokenString, nil
}

func (j JwtServiceImpl) GenerateRefreshToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.jwtRefreshExpiration) * time.Second)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.jwtRefreshSecret))
	if err != nil {
		j.logger.Error("Error generating refresh token", zap.Error(err))
		return "", fmt.Errorf("%w:%w", services.ErrGenerateToken, err)
	}

	return tokenString, nil
}

func (j JwtServiceImpl) ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.logger.Error("Unexpected signing method")
			return nil, fmt.Errorf("%w: %v", services.ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(j.jwtRefreshSecret), nil
	})

	if err != nil {
		j.logger.Error("Error parsing token")
		return nil, fmt.Errorf("%w:%w", services.ErrParseToken, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, services.ErrInvalidToken
}

func NewJwtServiceImpl(
	jwtSecret string,
	jwtRefreshSecret string,
	jwtExpiration int,
	jwtRefreshExpiration int,
	logger *zap.Logger,
) *JwtServiceImpl {
	return &JwtServiceImpl{
		jwtSecret,
		jwtRefreshSecret,
		jwtExpiration,
		jwtRefreshExpiration,
		logger,
	}
}
