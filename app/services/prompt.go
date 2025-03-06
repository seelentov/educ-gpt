package services

import "educ-gpt/models"

type PromptService interface {
	GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) (string, error)
	GetProblems(count int, topic string, theme string, userThemeStats *models.Theme, userAllStats []*models.Theme) (string, error)
	GetTheme(topic string, theme string, userStats *models.Theme, userAllStats []*models.Theme) (string, error)
	VerifyAnswer(problem string, answer string) (string, error)
}

type PromptProblemsRequest struct {
	Problems []string `json:"problems"`
}
type PromptThemeRequest struct {
	Text     string   `json:"text"`
	Problems []string `json:"problems"`
}
type PromptProblemRequest struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
