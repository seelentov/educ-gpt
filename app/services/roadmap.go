package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrGetEntities       = errors.New("cannot get entites")
	ErrUpdateEntity      = errors.New("cannot get entites")
	ErrCreateEntity      = errors.New("cannot create entity")
	ErrGetOrCreateEntity = errors.New("cannot get or create entity")
)

type RoadmapService interface {
	GetTopics(userID uint) ([]*models.Topic, error)
	GetTopic(userID uint, topicID uint) (*models.Topic, error)
	CreateThemes(theme []*models.Theme) error
	IncrementUserScore(userID uint, themeID uint, score uint) error
	UpdateUserResolvedProblems(userID uint, themeID uint, newProblem string) error
}
