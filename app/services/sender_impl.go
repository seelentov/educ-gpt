package services

import (
	"context"
	"crypto/tls"
	"educ-gpt/jobs/tasks"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type SenderServiceImpl struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string

	queueName string

	logger      *zap.Logger
	redisClient *redis.Client
}

func (s SenderServiceImpl) SendMessageByWorker(to, subject, body string) error {
	ctx := context.Background()
	task := tasks.MailTask{To: to, Subject: subject, Body: body}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		s.logger.Error("failed to marshal task", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrQueuedMail, err)
	}

	err = s.redisClient.RPush(ctx, "email_sender", taskJSON).Err()
	if err != nil {
		s.logger.Error("failed to push task to redis", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrQueuedMail, err)
	}

	return nil
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

func NewSenderServiceImpl(smtpHost string, smtpPort int, smtpUsername, smtpPassword, smtpFrom, queueName string, logger *zap.Logger, redisClient *redis.Client) SenderService {
	return &SenderServiceImpl{smtpHost, smtpPort, smtpUsername, smtpPassword, smtpFrom, queueName, logger, redisClient}
}
