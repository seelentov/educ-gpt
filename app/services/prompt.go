package services

import "educ-gpt/models"

type PromptService interface {
	GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) (string, error)
	GetTheme(topic string, theme string, userStats []*models.Theme) (string, error)
}

type PromptThemeRequest struct {
	Text     string
	Problems []string
}
