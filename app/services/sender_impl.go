package services

import (
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type SenderServiceImpl struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string

	logger *zap.Logger
}

func (s SenderServiceImpl) SendMessage(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		s.logger.Error("SendMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrSendMail, err)
	}

	return nil
}

func NewSenderServiceImpl(smtpHost string, smtpPort int, smtpUsername, smtpPassword, smtpFrom string, logger *zap.Logger) SenderService {
	return &SenderServiceImpl{smtpHost, smtpPort, smtpUsername, smtpPassword, smtpFrom, logger}
}
