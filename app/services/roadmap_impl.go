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

func (r RoadmapServiceImpl) GetProblem(problemID uint) (string, error) {
	var problem *models.Problem
	result := r.db.Model(&models.Problem{}).Where("id = ?", problemID).First(&problem)
	if result.Error != nil {
		r.logger.Error("Cant get problem", zap.Error(result.Error))
		return "", fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
	}

	return problem.Question, nil
}

func (r RoadmapServiceImpl) GetTheme(userID uint, themeID uint) (*models.Theme, error) {
	var theme *models.Theme
	r.db.Model(&models.Theme{}).Where("id = ?", themeID).First(&theme)

	if userID != 0 {
		var userTheme *models.UserTheme
		result := r.db.Model(&models.UserTheme{}).Where("user_id = ? AND theme_id = ?", userID, themeID).First(&userTheme)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				r.logger.Error("Cant get topic", zap.Error(result.Error))
				return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
			}
		} else {
			theme.Score = userTheme.Score
			theme.ResolvedProblems = userTheme.ResolvedProblems
		}
	}

	return theme, nil
}

func (r RoadmapServiceImpl) IncrementUserScoreAndAddAnswer(userID uint, themeID uint, newProblem string, score uint) error {
	userTheme := models.UserTheme{UserID: userID, ThemeID: themeID}

	result := r.db.Model(&models.UserTheme{}).FirstOrCreate(userTheme)
	if result.Error != nil {
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", ErrGetOrCreateEntity, result.Error)
	}

	userTheme.ResolvedProblems += "; " + newProblem
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

func (r RoadmapServiceImpl) GetTopics(userID uint, prThemes bool) ([]*models.Topic, error) {
	var topics []*models.Topic
	result := r.db.Model(&models.Topic{}).Preload("Themes")

	if prThemes {
		result = result.Preload("Themes")
	}

	result = result.Find(&topics)

	if result.Error != nil {
		r.logger.Error("Cant get topics", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
	}

	if userID != 0 && prThemes {
		for i := range topics {
			for j := range topics[i].Themes {
				var userTheme *models.UserTheme
				result := r.db.Model(&models.UserTheme{}).Where("user_id = ? AND theme_id = ?", userID, topics[i].Themes[j].ID).First(&userTheme)
				if result.Error != nil {
					if errors.Is(result.Error, gorm.ErrRecordNotFound) {
						continue
					}

					r.logger.Error("Cant get topic", zap.Error(result.Error))
					return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
				}

				topics[i].Themes[j].Score = userTheme.Score
				topics[i].Themes[j].ResolvedProblems = userTheme.ResolvedProblems
			}
		}
	}

	return topics, nil
}

func (r RoadmapServiceImpl) GetTopic(userID uint, topicID uint, prThemes bool) (*models.Topic, error) {
	var topic *models.Topic

	result := r.db.Model(&models.Topic{}).Where("id = ?", topicID)

	if prThemes {
		result = result.Preload("Themes")
	}

	result = result.First(&topic)

	if result.Error != nil {
		r.logger.Error("Cant get topic", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", ErrGetEntities, result.Error)
	}

	if userID != 0 && prThemes {
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
			topic.Themes[i].ResolvedProblems = userTheme.ResolvedProblems
		}
	}

	return topic, nil
}

func (r RoadmapServiceImpl) CreateProblems(problems []string) error {
	problemsStrs := make([]*models.Problem, len(problems))

	for i := range problems {
		problemsStrs[i] = &models.Problem{Question: problems[i]}
	}

	if err := r.db.Save(&problemsStrs).Error; err != nil {
		r.logger.Error("Cant create topics", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrCreateEntity, err)
	}

	return nil
}

func (r RoadmapServiceImpl) DeleteProblem(problemID uint) error {
	if err := r.db.Where("id = ?", problemID).Delete(&models.Problem{}).Error; err != nil {
		r.logger.Error("Cant create problem", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrDeleteEntities, err)
	}
	return nil
}

func NewRoadmapServiceImpl(db *gorm.DB, logger *zap.Logger) RoadmapService {
	return &RoadmapServiceImpl{db, logger}
}
