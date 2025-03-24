package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"fmt"
	"strings"
	"testing"
)

var (
	mailSrv services.MailService

	userIdMailSrv uint = 1
	nameMailSrv        = "test"
	keyMailSrv         = "key"
)

func TestCatInitMailService(t *testing.T) {
	mailSrv = dic.MailService()
}

func TestCanActivateMail(t *testing.T) {
	expSbj := "Активация аккаунта EDUC GPT"

	mail, err := mailSrv.ActivateMail(nameMailSrv, keyMailSrv)

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if !strings.Contains(mail.Body, nameMailSrv) {
		t.Error("Missing name")
		return
	}

	if !strings.Contains(mail.Body, keyMailSrv) {
		t.Error("Missing key")
		return
	}
}

func TestCanChangeEmailMail(t *testing.T) {
	expSbj := "Смена почтового ящика EDUC GPT"

	mail, err := mailSrv.ChangeEmailMail(userIdMailSrv, nameMailSrv, keyMailSrv)

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if !strings.Contains(mail.Body, fmt.Sprintf("/%v", userIdMailSrv)) {
		t.Error("Missing user id")
		return
	}

	if !strings.Contains(mail.Body, nameMailSrv) {
		t.Error("Missing name")
		return
	}

	if !strings.Contains(mail.Body, keyMailSrv) {
		t.Error("Missing key")
		return
	}
}

func TestCanResetMail(t *testing.T) {
	expSbj := "Восстановление аккаунта EDUC GPT"

	mail, err := mailSrv.ResetMail(userIdMailSrv, nameMailSrv, keyMailSrv)

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if !strings.Contains(mail.Body, fmt.Sprintf("/%v", userIdMailSrv)) {
		t.Error("Missing user id")
		return
	}

	if !strings.Contains(mail.Body, nameMailSrv) {
		t.Error("Missing name")
		return
	}

	if !strings.Contains(mail.Body, keyMailSrv) {
		t.Error("Missing key")
		return
	}
}
