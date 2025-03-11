package services

type Mail struct {
	Subject string
	Body    string
}

type MailService interface {
	ActivateMail(name, key string) (*Mail, error)
	ResetMail(userId uint, name, key string) (*Mail, error)
	ChangeEmailMail(userId uint, name, key string) (*Mail, error)
}
