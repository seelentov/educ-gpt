package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"testing"
)

var (
	tokenSrv services.TokenService

	keyTemp   string
	tokenData = "test_data"
)

func TestInitTokenService(t *testing.T) {
	tokenSrv = dic.TokenService()
}

func TestCanCreateToken(t *testing.T) {
	key, err := tokenSrv.Create(0, models.TypeResetPassword, tokenData)
	if err != nil {
		t.Error(err)
		return
	}

	keyTemp = key
}

func TestCanVerifyAndGetDataToken(t *testing.T) {
	data, err := tokenSrv.VerifyAndGetData(0, keyTemp, models.TypeResetPassword)
	if err != nil {
		t.Error(err)
		return
	}

	if data != tokenData {
		t.Errorf("Expected %s but got %s", tokenData, data)
		return
	}
}
