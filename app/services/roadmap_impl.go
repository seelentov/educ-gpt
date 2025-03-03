package services

import (
	"educ-gpt/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoadmapServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r RoadmapServiceImpl) UpdateUserResolvedProblems(userID uint, themeID uint, newProblem string) error {
	userTheme := models.UserTheme{UserID: userID, ThemeID: themeID}

	result := r.db.Model(&models.UserTheme{}).FirstOrCreate(userTheme)
	if result.Error != nil {
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", ErrGetOrCreateEntity, result.Error)
	}

	userTheme.ResolvedProblems += "; " + newProblem

	result = r.db.Save(&userTheme)
	if result.Error != nil {
		r.logger.Error("Cant update user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", ErrUpdateEntity, result.Error)
	}

	return nil

}

func (r RoadmapServiceImpl) IncrementUserScore(userID uint, themeID uint, score uint) error {
	userTheme := models.UserTheme{UserID: userID, ThemeID: themeID}

	result := r.db.Model(&models.UserTheme{}).FirstOrCreate(userTheme)
	if result.Error != nil {
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", ErrGetOrCreateEntity, result.Error)
	}

	userTheme.Score += score

	result = r.db.Save(&userTheme)
	if result.Error != nil {
		r.logger.Error("Cant update user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", ErrUpdateEntity, result.Error)
	}

	return nil
}

func (r RoadmapServiceImpl) CreateThemes(theme []*models.Theme) error {
	if err := r.db.Save(&theme).Error; err != nil {
		r.logger.Error("Cant create topics", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrCreateEntity, err)
	}
	return nil
}

func (r RoadmapServiceImpl) GetTopics(userID uint) ([]*models.Topic, error) {
	var topics []*models.Topic
	result := r.db.Model(&models.Topic{}).Preload("Themes").Find(&topics)

	if result.Error != nil {
		r.logger.Error("Cant get topics", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
	}

	for i := range topics {
		themeIds := make([]uint, len(topics[i].Themes))

		for j := range topics[i].Themes {
			themeIds[j] = topics[i].Themes[j].ID
		}

		var userThemes []*models.UserTheme
		result := r.db.Model(&models.UserTheme{}).Where("user_id = ? AND theme_id in (?)", userID, themeIds).Preload("Theme").Find(&userThemes)

		themes := make([]*models.Theme, len(userThemes))
		for j := range userThemes {
			themes[j] = userThemes[j].Theme
		}

		if result.Error != nil {
			r.logger.Error("Cant get topics", zap.Error(result.Error))
			return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
		}

		for _, theme := range themes {
			topics[i].Score += theme.Score
		}

	}

	return topics, nil
}

func (r RoadmapServiceImpl) GetTopic(userID uint, topicID uint) (*models.Topic, error) {
	var topic *models.Topic
	result := r.db.Model(&models.Topic{}).Where("id = ?", topicID).Preload("Themes").First(&topic)
	if result.Error != nil {
		r.logger.Error("Cant get topic", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
	}

	for i := range topic.Themes {
		var userTheme *models.UserTheme
		result := r.db.Model(&models.UserTheme{}).Where("user_id = ? AND theme_id = ?", userID, topic.Themes[i].ID).First(&userTheme)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				continue
			}

			r.logger.Error("Cant get topic", zap.Error(result.Error))
			return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
		}

		topic.Themes[i].Score = userTheme.Score
	}

	return topic, nil
}

func NewRoadmapServiceImpl(db *gorm.DB, logger *zap.Logger) RoadmapService {
	return &RoadmapServiceImpl{db, logger}
}
