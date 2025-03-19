package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"strconv"
	"testing"
)

var (
	srv = dic.UserService()

	activationToken string

	user = &models.User{
		Name:         "test_user",
		Email:        "test_user@test.com",
		Password:     "test_password",
		ChatGptToken: "test_token",
	}
)

func userFactory(i int) *models.User {
	iString := strconv.Itoa(i)

	return &models.User{
		Name:         "test_user" + iString,
		Email:        "test_user" + iString + "@test.com",
		Password:     "test_password",
		ChatGptToken: "test_token",
	}
}

func TestInitService(t *testing.T) {
	srv = dic.UserService()
}

func TestCreateUser(t *testing.T) {

	token, err := srv.Create(user)
	if err != nil {
		t.Error(err)
	}

	activationToken = token
}

func TestActivateUser(t *testing.T) {
	err := srv.Activate(activationToken)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUserById(t *testing.T) {
	u, err := srv.GetById(user.ID)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
	}
}

func TestGetUserByName(t *testing.T) {
	u, err := srv.GetByName(user.Name)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
	}
}

func TestGetUserByEmail(t *testing.T) {
	u, err := srv.GetByEmail(user.Email)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
	}
}

func TestGetUserByCredential(t *testing.T) {
	u, err := srv.GetByCredential(user.Name)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
	}

	u, err = srv.GetByCredential(user.Email)
	if err != nil {
		t.Error(err)
	}

	if u == nil {
		t.Errorf("Expected user but got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	updatesMap := make(map[string]interface{})

	newName := user.Name + "_updated"

	updatesMap["name"] = newName

	if err := srv.Update(user.ID, updatesMap); err != nil {
		t.Error(err)
	}

	u, err := srv.GetById(user.ID)

	if err != nil {
		t.Error(err)
	}

	if u.Name != newName {
		t.Errorf("Expected %s but got %s", newName, u.Name)
	}
}

func TestCanVerify(t *testing.T) {
	if err := srv.Verify(user.Password, user.Email); err != nil {
		t.Error(err)
	}
}

func TestCantVerifyIfWrong(t *testing.T) {
	if err := srv.Verify(user.Password+"_wrong", user.Email); err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestCanVerifyPassword(t *testing.T) {
	u, err := srv.GetById(user.ID)
	if err != nil {
		t.Error(err)
	}

	if err := srv.VerifyPassword(user.Password, u.Password); err != nil {
		t.Error(err)
	}
}

func TestCantVerifyPasswordIfWrong(t *testing.T) {
	if err := srv.VerifyPassword(user.Password, "wrong"); err != nil {
		t.Error(err)
	}
}

func TestCanChangePassword(t *testing.T) {
	newPassword := "test_new_password"
	if err := srv.ChangePassword(user.ID, newPassword); err != nil {
		t.Error(err)
	}

	if err := srv.Verify(newPassword, user.Email); err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestCanClearNonActivatedUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		if _, err := srv.Create(userFactory(i)); err != nil {
			t.Error(err)
		}
	}

	if err := srv.ClearNonActivatedUsers(); err != nil {
		t.Error(err)
	}

	for i := 1; i < 11; i++ {
		if u, _ := srv.GetById(user.ID + uint(i)); u != nil {
			t.Errorf("Expected nil but got user with id %v", i)
		}
	}
}
