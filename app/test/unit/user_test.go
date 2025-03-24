package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"testing"
	"time"
)

var (
	userSrv services.UserService

	activationToken string

	user = &models.User{
		Name:     "test_user",
		Email:    "test_user@test.com",
		Password: "test_password",
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

func TestCanClearNonActivatedUsers(t *testing.T) {

	nonActivatedUser := &models.User{
		Name:      "non_activated_user",
		Email:     "non_activated_user@test.com",
		Password:  "non_activated_password",
		CreatedAt: time.Now().Add(-3 * time.Hour),
	}

	_, err := userSrv.Create(nonActivatedUser)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	createdUser, err := userSrv.GetByEmail(nonActivatedUser.Email)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if createdUser == nil {
		t.Error("Expected user to be created, but got nil")
		return
	}

	err = userSrv.ClearNonActivatedUsers()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	deletedUser, err := userSrv.GetByEmail(nonActivatedUser.Email)
	if err == nil {
		t.Error("Expected error when getting deleted user, but got nil")
		return
	}

	if deletedUser != nil {
		t.Error("Expected user to be deleted, but got a user")
		return
	}
}
