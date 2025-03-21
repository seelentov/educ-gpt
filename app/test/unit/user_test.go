package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"testing"
)

var (
	userSrv services.UserService

	activationToken string

	user = &models.User{
		Name:         "test_user",
		Email:        "test_user@test.com",
		Password:     "test_password",
		ChatGptToken: "test_token",
	}
)

func TestInitUserService(t *testing.T) {
	userSrv = dic.UserService()
}

func TestCreateUser(t *testing.T) {
	tempPass := user.Password
	token, err := userSrv.Create(user)
	user.Password = tempPass
	if err != nil {
		t.Error(err)
		return
	}

	activationToken = token
}

func TestActivateUser(t *testing.T) {
	err := userSrv.Activate(activationToken)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUserById(t *testing.T) {
	u, err := userSrv.GetById(user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
		return
	}
}

func TestGetUserByName(t *testing.T) {
	u, err := userSrv.GetByName(user.Name)
	if err != nil {
		t.Error(err)
		return
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
		return
	}
}

func TestGetUserByEmail(t *testing.T) {
	u, err := userSrv.GetByEmail(user.Email)
	if err != nil {
		t.Error(err)
		return
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
		return
	}
}

func TestGetUserByCredential(t *testing.T) {
	u, err := userSrv.GetByCredential(user.Name)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
		return
	}

	u, err = userSrv.GetByCredential(user.Email)
	if err != nil {
		t.Error(err)
		return
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
		return
	}
}

func TestUpdateUser(t *testing.T) {
	updatesMap := make(map[string]interface{})

	newName := user.Name + "_updated"

	updatesMap["name"] = newName

	if err := userSrv.Update(user.ID, updatesMap); err != nil {
		t.Error(err)
		return
	}

	u, err := userSrv.GetById(user.ID)

	if err != nil {
		t.Error(err)
		return
	}

	if u.Name != newName {
		t.Errorf("Expected %s but got %s", newName, u.Name)
		return
	}
}

func TestCanVerifyPassword(t *testing.T) {
	u, err := userSrv.GetById(user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if err := userSrv.VerifyPassword(user.Password, u.Password); err != nil {
		t.Error(err)
	}
}

func TestCantVerifyPasswordIfWrong(t *testing.T) {
	u, err := userSrv.GetById(user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if err := userSrv.VerifyPassword("wrong", u.Password); err == nil {
		t.Error("Expected error but got nil")
		return
	}
}

func TestCanChangePassword(t *testing.T) {
	newPassword := "test_new_password"
	if err := userSrv.ChangePassword(user.ID, newPassword); err != nil {
		t.Error(err)
		return
	}

	u, err := userSrv.GetById(user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if err := userSrv.VerifyPassword(newPassword, u.Password); err != nil {
		t.Error(err)
		return
	}
}
