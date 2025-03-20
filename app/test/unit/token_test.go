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

func TestCanInit(t *testing.T) {
	tokenSrv = dic.TokenService()
}

func TestCanCreate(t *testing.T) {
	key, err := tokenSrv.Create(0, models.TypeResetPassword, tokenData)
	if err != nil {
		t.Error(err)
	}

	keyTemp = key
}

func TestCanVerifyAndGetData(t *testing.T) {
	data, err := tokenSrv.VerifyAndGetData(0, keyTemp, models.TypeResetPassword)
	if err != nil {
		t.Error(err)
	}

	if data != tokenData {
		t.Errorf("Expected %s but got %s", tokenData, data)
	}
}
