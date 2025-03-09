package services

import "errors"

var (
	ErrParsingTemplate = errors.New("error parsing template")
	ErrGenBody         = errors.New("error generating body")
)

type Mail struct {
	Subject string
	Body    string
}

type MailService interface {
	ActivateMail(name, activationKey string) (*Mail, error)
	ResetMail(name, activationKey string) (*Mail, error)
}
