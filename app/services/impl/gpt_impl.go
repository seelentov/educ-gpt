package impl

import (
	"bytes"
	"educ-gpt/models"
	"educ-gpt/services"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type GptServiceImpl struct {
	logger *zap.Logger
}

func (g GptServiceImpl) GetAnswer(token string, model string, dialog []*models.DialogItem, target interface{}) error {
	dialogStrings := make([]string, len(dialog))

	for i := range dialog {
		role := "assistant"

		if dialog[i].IsUser {
			role = "user"
		}

		dialogStrings[i] = fmt.Sprintf(`{"role": "%s", "content": "%s"}`, role, dialog[i].Text)
	}

	body := fmt.Sprintf(`{"model": "%s","messages": [%s],"temperature": 0.1}`, model, strings.Join(dialogStrings, ","))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(body)))
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrRequestFailed, err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	g.logger.Debug("request body", zap.String("body", body))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrRequestFailed, err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrParseFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%v:%s", services.ErrAIRequestFailed, resp.StatusCode, bodyBytes)
	}

	tempTarget := &services.GptResponse{}

	err = json.Unmarshal(bodyBytes, tempTarget)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrParseResFailed, err)
	}

	msg := tempTarget.Choices[len(tempTarget.Choices)-1].Message.Content

	msg = strings.TrimPrefix(msg, "```json")
	msg = strings.TrimSuffix(msg, "```")

	if reflect.TypeOf(target).String() == "*string" {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(msg))
		return nil
	}

	g.logger.Debug("received response", zap.String("message", msg))

	err = json.Unmarshal([]byte(msg), &target)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrParseFailed, err)
	}

	return nil
}

func NewGptService(logger *zap.Logger) services.GptService {
	return &GptServiceImpl{logger}
}
