package services

type Mail struct {
	Subject string
	Body    string
}

type MailService interface {
	ActivateMail(name, key string) (*Mail, error)
	ResetMail(name, key string) (*Mail, error)
	ChangeEmailMail(name, key string) (*Mail, error)
}
