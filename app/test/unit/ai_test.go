package unit

import (
	"educ-gpt/config/dic"
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
	var res string

	answer := "Ответь одной цифрой на вопрос: Сколько будет 1+1?"

	err := aiSrv.GetAnswer(
		os.Getenv("OPENROUTER_TOKEN"), os.Getenv("OPENROUTER_MODEL"),
		[]*models.DialogItem{{Text: answer, IsUser: true}},
		&res,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res != "2" {
		t.Errorf("Expected 2, but got %s", res)
	}
}
