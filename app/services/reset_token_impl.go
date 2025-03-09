package services

import (
	"educ-gpt/models"
	"educ-gpt/utils/securityutils"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ResetTokenServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r ResetTokenServiceImpl) Create(userID uint) (string, error) {
	token := &models.ResetToken{}

	token.UserID = userID

	key := securityutils.GenerateKey(200)

	token.Key = key

	if err := r.db.Create(&token).Error; err != nil {
		r.logger.Error("Error creating token", zap.Error(err))
		return "", fmt.Errorf("%w: %w", ErrCreateResetToken, err)
	}

	return key, nil
}

func (r ResetTokenServiceImpl) Verify(userID uint, key string) error {
	var tokens []*models.ResetToken

	if err := r.db.Where("user_id = ? AND key = ?", userID, key).Find(&tokens).Error; err != nil {
		r.logger.Error("Error verifying token", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrVerifyResetToken, err)
	}

	if len(tokens) == 0 {
		r.logger.Error("Error verifying token: token not found")
		return fmt.Errorf("%w: %w", ErrVerifyResetToken, gorm.ErrRecordNotFound)
	}

	ago := time.Now().Add(-2 * time.Hour)
	ok := false
	for _, token := range tokens {
		if token.CreatedAt.After(ago) {
			ok = true
			break
		}
	}

	if !ok {
		r.logger.Error("Error verifying token: token expired")
		return fmt.Errorf("%w: %w", ErrVerifyResetToken)
	}

	if err := r.db.Where("user_id = ?", userID).Delete(&models.ResetToken{}).Error; err != nil {
		r.logger.Error("Error clear tokens", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrClearResetTokens, err)
	}

	return nil
}

func (r ResetTokenServiceImpl) Clear() error {
	ago := time.Now().Add(-2 * time.Hour)

	result := r.db.Where("created_at < ?", ago).Delete(&models.ResetToken{})

	if result.Error != nil {
		r.logger.Error("Error deleting old reset tokens", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", ErrVerifyResetToken, result.Error)
	}

	return nil
}

func NewResetTokenServiceImpl(db *gorm.DB, logger *zap.Logger) ResetTokenService {
	return &ResetTokenServiceImpl{db, logger}
}
