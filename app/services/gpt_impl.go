package services

import (
	"bytes"
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

func (g GptServiceImpl) GetAnswer(token string, model string, dialog []*DialogItem, target interface{}) error {
	dialogStrings := make([]string, len(dialog))

	for i := range dialog {
		role := "assistant"

		if dialog[i].IsUser {
			role = "user"
		}

		dialogStrings[i] = fmt.Sprintf("{\"role\": \"%s\", \"content\": \"%s\"}", role, dialog[i].Text)
	}

	body := fmt.Sprintf(`{
     "model": "%s",
     "messages": [%s],
     "temperature": 0.7
   	}`, model, strings.Join(dialogStrings, ","))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(body)))
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrRequestFailed, err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrRequestFailed, err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrParseFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w:%v:%s", ErrRequestFailed, resp.StatusCode, bodyBytes)
	}

	tempTarget := &GptResponse{}

	err = json.Unmarshal(bodyBytes, tempTarget)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrParseFailed, err)
	}

	msg := tempTarget.Choices[len(tempTarget.Choices)-1].Message.Content

	if reflect.TypeOf(target).String() == "*string" {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(msg))
		return nil
	}

	err = json.Unmarshal([]byte(msg), &target)
	if err != nil {
		g.logger.Error("failed to send request", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrParseFailed, err)
	}

	return nil
}

func NewGptService(logger *zap.Logger) GptService {
	return &GptServiceImpl{logger}
}
