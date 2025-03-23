package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"os"
	"testing"
)

var (
	aiSrv services.AIService
)

func TestCanInitAIService(t *testing.T) {
	aiSrv = dic.AIService()
}

func TestCanGetAnswer(t *testing.T) {
	var res dtos.ResultResponse

	answer := "Ответь одной цифрой на вопрос: Сколько будет 1+1?"

	err := aiSrv.GetAnswer(
		os.Getenv("ADMIN_CHAT_GPT_TOKEN"), os.Getenv("ADMIN_CHAT_GPT_MODEL"),
		[]*models.DialogItem{{Text: answer, IsUser: true}},
		&res,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Result != "2" {
		t.Errorf("Expected 2, but got %s", res)
	}
}
