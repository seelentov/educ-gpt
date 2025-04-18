package services

import "errors"

var (
	ErrGetEntities       = errors.New("cannot get entites")
	ErrDeleteEntities    = errors.New("cannot get delete")
	ErrCreateEntity      = errors.New("cannot create entity")
	ErrGetOrCreateEntity = errors.New("cannot get or create entity")
)

type RoadmapService interface {
	TopicService
	ProblemService
	ThemeService
}
