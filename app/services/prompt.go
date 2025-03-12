package services

import "educ-gpt/models"

type PromptService interface {
	GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) string
	GetProblems(count int, topic string, theme string, userThemeStats *models.Theme, userAllStats []*models.Theme) string
	GetTheme(topic string, theme string, userStats *models.Theme, userAllStats []*models.Theme) string
	VerifyAnswer(problem string, answer string) string
	CompileCode(code string) string
}

type PromptProblemsResponse struct {
	Problems []string `json:"problems"`
}
type PromptThemeResponse struct {
	Text     string   `json:"text"`
	Problems []string `json:"problems"`
}
type PromptProblemRequest struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
