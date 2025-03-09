package services

import "errors"

var (
	ErrSendMail = errors.New("error send mail")
)

type SenderService interface {
	SendMessage(to, subject, body string) error
}
