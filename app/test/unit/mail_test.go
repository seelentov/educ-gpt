package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"os"
	"testing"
)

var (
	mailSrv services.MailService
)

func TestCatInitMailService(t *testing.T) {
	mailSrv = dic.MailService()
}

func TestCanActivateMail(t *testing.T) {
	expSbj := "Активация аккаунта EDUC GPT"
	expBody := `<html><body><h1>Привет, test!</h1><p>Спасибо за регистрацию на EDUC GPT. Пожалуйста, активируйте ваш аккаунт, перейдя по ссылке ниже:</p><a href="` + os.Getenv("PROTOCOL") + "://" + os.Getenv("HOST") + `/activate/key">Активировать аккаунт</a><p><small>Если вы не регистрировали аккаунт на EDUC GPT, то проигнорируйте это письмо</small></p></body></html>`

	mail, err := mailSrv.ActivateMail("test", "key")

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if expBody != mail.Body {
		t.Errorf("Expected %s but got %s", expBody, mail.Body)
		return
	}
}

func TestCanChangeEmailMail(t *testing.T) {
	expSbj := "Смена почтового ящика EDUC GPT"
	expBody := `<html><body><h1>Привет, test!</h1><p>Для смены почтового ящика на своем аккаунте EDUC GPT на этот перейдите по ссылке ниже:</p><a href="` + os.Getenv("PROTOCOL") + "://" + os.Getenv("HOST") + `/change_email/key/1">Сменить почту</a><p>Ссылка будет активна 2 часа</p><p><small>Если вы не пытались сменить почтовый ящик своего аккаунта на EDUC GPT, то проигнорируйте это письмо</small></p></body></html>`

	mail, err := mailSrv.ChangeEmailMail(1, "test", "key")

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if expBody != mail.Body {
		t.Errorf("Expected %s but got %s", expBody, mail.Body)
		return
	}
}

func TestCanResetMail(t *testing.T) {
	expSbj := "Восстановление аккаунта EDUC GPT"
	expBody := `<html><body><h1>Привет, test!</h1><p>Для смены пароля на своем аккаунте EDUC GPT перейдите по ссылке ниже:</p><a href="` + os.Getenv("PROTOCOL") + "://" + os.Getenv("HOST") + `/reset/key/1">Сменить пароль</a><p>Ссылка будет активна 2 часа</p><p><small>Если вы не пытались восстановить пароль своего аккаунта на EDUC GPT, то проигнорируйте это письмо</small></p></body></html>`

	mail, err := mailSrv.ResetMail(1, "test", "key")

	if err != nil {
		t.Error(err)
		return
	}

	if expSbj != mail.Subject {
		t.Errorf("Expected %s but got %s", expSbj, mail.Subject)
		return
	}

	if expBody != mail.Body {
		t.Errorf("Expected %s but got %s", expBody, mail.Body)
		return
	}
}
