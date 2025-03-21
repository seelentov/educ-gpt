package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"os"
	"testing"
)

var (
	roleService services.RoleService
	userId      uint = 1
)

func TestCanInitRoleService(t *testing.T) {
	roleService = dic.RoleService()
}

func TestGetRoleByName(t *testing.T) {
	role, err := roleService.GetRoleByName("USER")
	if err != nil {
		t.Error(err)
		return
	}

	if role == nil {
		t.Error("Role not found")
		return
	}

	if role.Name != "USER" {
		t.Errorf("Expected role name to be 'USER', got %v", role.Name)
		return
	}
}

func TestCanGetRolesByUserId(t *testing.T) {
	roles, err := roleService.GetRolesByUserId(userId)
	if err != nil {
		t.Error(err)
		return
	}

	exist := false

	for _, role := range roles {
		if role.Name == "ADMIN" {
			exist = true
			break
		}
	}

	if !exist {
		t.Error("Role ADMIN by UserID 1 not found")
		return
	}
}

func TestGetUsersByName(t *testing.T) {
	users, err := roleService.GetUsersByName("USER")
	if err != nil {
		t.Error(err)
		return
	}

	if len(users) == 0 {
		t.Error("Users by role 'USER' not found")
		return
	}
}

func TestGetDefaultRole(t *testing.T) {
	role, err := roleService.GetDefaultRole()
	if err != nil {
		t.Error(err)
		return
	}

	if role == nil {
		t.Error("Default role not found")
		return
	}

	if role.Name != os.Getenv("AUTH_DEFAULT_ROLE") {
		t.Errorf("Expected role name to be %s, got %s", os.Getenv("AUTH_DEFAULT_ROLE"), role.Name)
	}
}
