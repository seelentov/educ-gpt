package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrGetEntities       = errors.New("cannot get entites")
	ErrDeleteEntities    = errors.New("cannot get delete")
	ErrUpdateEntity      = errors.New("cannot get entites")
	ErrCreateEntity      = errors.New("cannot create entity")
	ErrGetOrCreateEntity = errors.New("cannot get or create entity")
)

type RoadmapService interface {
	IncrementUserScoreAndAddAnswer(userID uint, themeID uint, newProblem string, score uint) error
	CreateThemes(theme []*models.Theme) error
	GetTopics(userID uint, prThemes bool) ([]*models.Topic, error)
	GetTopic(userID uint, topicID uint, prThemes bool) (*models.Topic, error)
	GetTheme(userID uint, themeID uint) (*models.Theme, error)
	GetProblem(problemID uint) (*models.Problem, error)
	CreateProblems(problems []string, themeID uint) ([]*models.Problem, error)
	DeleteProblem(problemID uint) error
}
