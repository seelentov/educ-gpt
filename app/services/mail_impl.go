package services

import (
	"fmt"
	"go.uber.org/zap"
)

type MailServiceImpl struct {
	protocol       string
	host           string
	activationLink string
	resetLink      string
	logger         *zap.Logger
}

func (s *MailServiceImpl) ResetMail(name, key string) (*Mail, error) {
	link := s.protocol + "://" + s.host + "/" + s.resetLink + "/" + key

	emailTemplate := fmt.Sprintf(`<html><body><h1>Привет, %s!</h1><p>Для смены пароля на своем аккаунте EDUC GPT перейдите по ссылке ниже:</p><a href="%s">Сменить пароль</a><p>Ссылка будет активна 2 часа</p><p><small>Если вы не пытались восстановить пароль своего аккаунта на EDUC GPT, то проигнорируйте это письмо</small></p></body></html>`, name, link)

	mail := &Mail{
		Subject: "Активация аккаунта EDUC GPT",
		Body:    emailTemplate,
	}

	return mail, nil
}

func (s *MailServiceImpl) ActivateMail(name, activationKey string) (*Mail, error) {
	activationLink := s.protocol + "://" + s.host + "/" + s.activationLink + "/" + activationKey

	emailTemplate := fmt.Sprintf(`<html><body><h1>Привет, %s!</h1><p>Спасибо за регистрацию на EDUC GPT. Пожалуйста, активируйте ваш аккаунт, перейдя по ссылке ниже:</p><a href="%s">Активировать аккаунт</a><p><small>Если вы не регистрировали аккаунт на EDUC GPT, то проигнорируйте это письмо</small></p></body></html>`, name, activationLink)

	mail := &Mail{
		Subject: "Восстановление аккаунта EDUC GPT",
		Body:    emailTemplate,
	}

	return mail, nil
}

func NewMailServiceImpl(protocol, host, activationLink, resetLink string, logger *zap.Logger) MailService {
	return &MailServiceImpl{protocol, host, activationLink, resetLink, logger}
}
