package services

import "educ-gpt/models"

type PromptService interface {
	GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) (string, error)
	GetTheme(topic string, theme string, userStats *models.Theme) (string, error)
	VerifyAnswer(problem string, answer string) (string, error)
}

type PromptThemeRequest struct {
	Text     string   `json:"text"`
	Problems []string `json:"problems"`
}

type PromptProblemRequest struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
