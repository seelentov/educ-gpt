package services

import "educ-gpt/models"

type ThemeService interface {
	IncrementUserScoreAndAddAnswer(userID uint, problemID uint, score uint) error
	CreateThemes(theme []*models.Theme) error
	GetTheme(userID uint, themeID uint, prTopic bool) (*models.Theme, error)
}
