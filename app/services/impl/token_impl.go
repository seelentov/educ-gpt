package impl

import (
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/securityutils"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type TokenServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r TokenServiceImpl) VerifyAndGetData(userID uint, key string, t models.TokenType) (string, error) {
	var token *models.Token

	if err := r.db.Where("user_id = ? AND key = ? AND type = ? AND created_at > ?", userID, key, t, time.Now().Add(-2*time.Hour)).First(&token).Error; err != nil {
		r.logger.Error("Error verifying token", zap.Error(err))
		return "", fmt.Errorf("%w: %w", services.ErrVerifyToken, err)
	}

	data := token.Data

	if err := r.db.Where("user_id = ? AND type = ?", userID, t).Delete(&models.Token{}).Error; err != nil {
		r.logger.Error("Error clear tokens", zap.Error(err))
		return "", fmt.Errorf("%w: %w", services.ErrClearTokens, err)
	}

	return data, nil
}

func (r TokenServiceImpl) Create(userID uint, t models.TokenType, data string) (string, error) {
	token := &models.Token{}

	token.UserID = userID

	key := securityutils.GenerateKey(200)

	token.Key = key
	token.Type = t
	token.Data = data

	if err := r.db.Create(&token).Error; err != nil {
		r.logger.Error("Error creating token", zap.Error(err))
		return "", fmt.Errorf("%w: %w", services.ErrCreateToken, err)
	}

	return key, nil
}

func (r TokenServiceImpl) Verify(userID uint, key string, t models.TokenType) error {
	_, err := r.VerifyAndGetData(userID, key, t)

	if err != nil {
		r.logger.Error("Error verifying token", zap.Error(err))
		return fmt.Errorf("%w: %w", services.ErrVerifyToken, err)
	}

	return nil
}

func (r TokenServiceImpl) Clear() error {
	threshold := time.Now().Add(-2 * time.Hour)

	result := r.db.Where("created_at < ?", threshold).Delete(&models.Token{})

	if result.Error != nil {
		r.logger.Error("Error deleting old reset tokens", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", services.ErrVerifyToken, result.Error)
	}

	return nil
}

func NewTokenServiceImpl(db *gorm.DB, logger *zap.Logger) services.TokenService {
	return &TokenServiceImpl{db, logger}
}
