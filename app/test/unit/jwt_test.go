package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"testing"
)

var (
	jwtSrv services.JwtService

	tokenTemp        string
	refreshTokenTemp string

	userIDJwtService uint = 1
)

func TestCanInitJwtService(t *testing.T) {
	jwtSrv = dic.JwtService()
}

func TestCanGenerateToken(t *testing.T) {
	token, err := jwtSrv.GenerateToken(userIDJwtService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenTemp = token
}

func TestCanGenerateRefreshToken(t *testing.T) {
	token, err := jwtSrv.GenerateRefreshToken(userIDJwtService)

	if err != nil {
		t.Error(err)
		return
	}

	refreshTokenTemp = token
}

func TestCanValidateToken(t *testing.T) {
	claims, err := jwtSrv.ValidateToken(tokenTemp)
	if err != nil {
		t.Error(err)
		return
	}

	userID := uint(claims["user_id"].(float64))

	if userID != userIDJwtService {
		t.Errorf("Expected user id to be %d, got %d", userIDJwtService, userID)
		return
	}
}

func TestCanValidateRefreshToken(t *testing.T) {
	claims, err := jwtSrv.ValidateRefreshToken(refreshTokenTemp)
	if err != nil {
		t.Error(err)
		return
	}

	userID := uint(claims["user_id"].(float64))

	if userID != userIDJwtService {
		t.Errorf("Expected user id to be %d, got %d", userIDJwtService, userID)
		return
	}
}
