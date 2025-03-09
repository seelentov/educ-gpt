package services

import "educ-gpt/models"

type ProblemService interface {
	GetProblem(problemID uint) (*models.Problem, error)
	CreateProblems(problems []string, themeID uint) ([]*models.Problem, error)
	DeleteProblem(problemID uint) error
	ClearProblems() error
}
