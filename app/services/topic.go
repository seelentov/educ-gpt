package services

import (
	"educ-gpt/models"
)

type TopicService interface {
	GetTopics(userID uint, prThemes bool) ([]*models.Topic, error)
	GetTopic(userID uint, topicID uint, prThemes bool) (*models.Topic, error)
}
