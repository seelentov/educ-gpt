package services

import "errors"

var (
	ErrSendMail   = errors.New("error send mail")
	ErrQueuedMail = errors.New("error queued mail")
)

type SenderService interface {
	SendMessage(to, subject, body string) error
	SendMessageByWorker(to, subject, body string) error
}
