package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"os"
	"testing"
)

var (
	senderSrv services.SenderService
)

func TestCanInitSenderService(t *testing.T) {
	senderSrv = dic.SenderService()
}

func TestCanSendMessage(t *testing.T) {
	err := senderSrv.SendMessage(
		os.Getenv("ADMIN_EMAIL"),
		"This is test",
		os.Getenv("<b>testing</b>"),
	)
	if err != nil {
		t.Error(err)
		return
	}
}
