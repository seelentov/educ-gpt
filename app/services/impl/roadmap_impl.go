package impl

import (
	"educ-gpt/models"
	"educ-gpt/services"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoadmapServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r RoadmapServiceImpl) GetTheme(userID uint, themeID uint, prTopic bool) (*models.Theme, error) {
	var theme models.Theme
	query := r.db.Model(&models.Theme{}).Where("id = ?", themeID)

	if prTopic {
		query = query.Preload("Topic")
	}

	if userID != 0 {
		query = query.Preload("UserThemes", "user_id = ?", userID)
	}

	result := query.First(&theme)

	if result.Error != nil {
		r.logger.Error("Cant get theme", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", services.ErrGetEntities, result.Error)
	}

	if userID != 0 && len(theme.UserThemes) > 0 {
		theme.Score = theme.UserThemes[0].Score
		theme.ResolvedProblems = theme.UserThemes[0].ResolvedProblems
	}

	return &theme, nil
}

func (r RoadmapServiceImpl) IncrementUserScoreAndAddAnswer(userID uint, problemID uint, score uint) error {
	problem, err := r.GetProblem(problemID)
	if err != nil {
		r.logger.Error("Cant get or create user_theme", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrGetOrCreateEntity, err)
	}

	userTheme := &models.UserTheme{UserID: userID, ThemeID: problem.ThemeID}

	result := r.db.Model(models.UserTheme{}).Where("theme_id = ? AND user_id = ?", problem.ThemeID, userID).FirstOrCreate(&userTheme)
	if result.Error != nil {
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrGetOrCreateEntity, result.Error)
	}

	userTheme.ResolvedProblems += "; " + problem.Question
	userTheme.Score += score

	tx := r.db.Begin()

	result = tx.Save(&userTheme)
	if result.Error != nil {
		tx.Rollback()
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrGetOrCreateEntity, result.Error)
	}

	result = tx.Delete(&models.Problem{}, problemID)
	if result.Error != nil {
		tx.Rollback()
		r.logger.Error("Cant get or create user_theme", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrGetOrCreateEntity, result.Error)
	}

	tx.Commit()

	return nil
}

func (r RoadmapServiceImpl) CreateThemes(theme []*models.Theme) error {
	if err := r.db.Save(&theme).Error; err != nil {
		r.logger.Error("Cant create topics", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrCreateEntity, err)
	}
	return nil
}

func (r RoadmapServiceImpl) GetTopics(userID uint, prThemes bool) ([]*models.Topic, error) {
	var topics []*models.Topic
	query := r.db.Model(&models.Topic{})

	if prThemes || userID != 0 {
		query = query.Preload("Themes.UserThemes", "user_id = ?", userID)
	}

	result := query.Find(&topics)

	if result.Error != nil {
		r.logger.Error("Cant get topics", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", services.ErrGetEntities, result.Error)
	}

	if userID != 0 {
		for i := range topics {
			for j := range topics[i].Themes {
				if len(topics[i].Themes[j].UserThemes) > 0 {
					topics[i].Score += topics[i].Themes[j].UserThemes[0].Score
				}
			}
		}
	}

	if !prThemes && userID != 0 {
		for i := range topics {
			topics[i].Themes = nil
		}
	}

	return topics, nil
}

func (r RoadmapServiceImpl) GetTopic(userID uint, topicID uint, prThemes bool) (*models.Topic, error) {
	var topic models.Topic
	query := r.db.Model(&models.Topic{}).Where("id = ?", topicID)

	if prThemes {
		query = query.Preload("Themes.UserThemes", "user_id = ?", userID)
	}

	result := query.First(&topic)

	if result.Error != nil {
		r.logger.Error("Cant get topic", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", services.ErrGetEntities, result.Error)
	}

	if userID != 0 && prThemes {
		for i := range topic.Themes {
			if len(topic.Themes[i].UserThemes) > 0 {
				topic.Themes[i].Score = topic.Themes[i].UserThemes[0].Score
				topic.Themes[i].ResolvedProblems = topic.Themes[i].UserThemes[0].ResolvedProblems
			}
		}
	}

	return &topic, nil
}

func (r RoadmapServiceImpl) CreateProblems(problems []*models.Problem) ([]*models.Problem, error) {
	if err := r.db.Save(&problems).Error; err != nil {
		r.logger.Error("Cant create problems", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", services.ErrCreateEntity, err)
	}

	return problems, nil
}

func (r RoadmapServiceImpl) DeleteProblem(problemID uint) error {
	if err := r.db.Where("id = ?", problemID).Delete(&models.Problem{}).Error; err != nil {
		r.logger.Error("Cant create problem", zap.Error(err))
		return fmt.Errorf("%w:%w", services.ErrDeleteEntities, err)
	}
	return nil
}

func (r RoadmapServiceImpl) ClearProblems() error {
	threshold := time.Now().Add(-24 * time.Hour)

	result := r.db.Where("created_at < ?", threshold).Delete(&models.Problem{})
	if result.Error != nil {
		r.logger.Error("Cant remove problem", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrDeleteEntities, result.Error)
	}
	return nil
}

func (r RoadmapServiceImpl) ClearThemes() error {
	threshold := time.Now().Add(-24 * time.Hour)

	result := r.db.
		Where("created_at < ? AND id NOT IN (SELECT theme_id FROM user_themes WHERE score > 0)", threshold).
		Delete(&models.Theme{})

	if result.Error != nil {
		r.logger.Error("Cant remove themes", zap.Error(result.Error))
		return fmt.Errorf("%w:%w", services.ErrDeleteEntities, result.Error)
	}

	return nil
}

func (r RoadmapServiceImpl) GetProblem(problemID uint) (*models.Problem, error) {
	var problem models.Problem
	result := r.db.Model(&models.Problem{}).Where("id = ?", problemID).First(&problem)
	if result.Error != nil {
		r.logger.Error("Cant get problem", zap.Error(result.Error))
		return nil, fmt.Errorf("%w:%w", services.ErrGetEntities, result.Error)
	}

	return &problem, nil
}

func NewRoadmapServiceImpl(db *gorm.DB, logger *zap.Logger) services.RoadmapService {
	return &RoadmapServiceImpl{db, logger}
}
