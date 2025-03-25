package impl

import (
	"bytes"
	"educ-gpt/models"
	"educ-gpt/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

type OpenRouterServiceImpl struct {
	logger *zap.Logger
}

func (g OpenRouterServiceImpl) GetAnswer(token string, model string, dialog []*models.DialogItem, target interface{}) error {
	dialogStrings := make([]string, len(dialog))

	for i := range dialog {
		role := "assistant"

		if dialog[i].IsUser {
			role = "user"
		}

		dialogStrings[i] = fmt.Sprintf(`{"role": "%s", "content": "%s"}`, role, strings.ReplaceAll(dialog[i].Text, `\n`, `\\n`))
	}

	body := fmt.Sprintf(`{"model": "%s","messages": [%s],"temperature": 0.1}`, model, strings.Join(dialogStrings, ","))

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer([]byte(body)))
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
		g.logger.Error("failed to send request", zap.String("resp", string(bodyBytes)))
		return fmt.Errorf("%w:%v:%s", services.ErrAIRequestFailed, resp.StatusCode, bodyBytes)
	}

	tempTarget := &services.AIResponse{}

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

func NewOpenRouterServiceImpl(logger *zap.Logger) services.AIService {
	return &OpenRouterServiceImpl{logger}
}
