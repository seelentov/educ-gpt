package services

import (
	"educ-gpt/models"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DialogServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (d DialogServiceImpl) RemoveDialog(dialogID uint) error {
	if err := d.db.Model(&models.Dialog{}).Where("id = ?", dialogID).Delete(&models.Dialog{}).Error; err != nil {
		d.logger.Error("RemoveDialog error", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrDeleteDialog, err)
	}

	return nil
}

func (d DialogServiceImpl) GetDialogsByUserID(userID uint) ([]*models.Dialog, error) {
	var dialog []*models.Dialog
	if err := d.db.Model(&models.Dialog{}).
		Where("user_id = ?", userID).
		Preload("DialogItems", func(tx *gorm.DB) *gorm.DB {
			return tx.Joins("INNER JOIN (SELECT dialog_id, MAX(created_at) AS max_created_at FROM dialog_items GROUP BY dialog_id) AS latest ON dialog_items.dialog_id = latest.dialog_id AND dialog_items.created_at = latest.max_created_at")
		}).
		Find(&dialog).Error; err != nil {
		d.logger.Error("GetDialogs error", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrGetDialogs, err)
	}

	return dialog, nil
}

func (d DialogServiceImpl) GetDialog(dialogID uint) (*models.Dialog, error) {
	var dialog models.Dialog
	if err := d.db.Model(&models.Dialog{}).Where("id = ?", dialogID).Preload("DialogItems").First(&dialog).Error; err != nil {
		d.logger.Error("GetDialog error", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrDialogNotFound, err)
	}

	return &dialog, nil
}

func (d DialogServiceImpl) AddDialogItem(dialogItem *models.DialogItem) error {
	if err := d.db.Create(dialogItem).Error; err != nil {
		d.logger.Error("AddDialogItem error", zap.Error(err))
		return fmt.Errorf("%w:%w", ErrAddDialogItem, err)
	}

	return nil
}

func (d DialogServiceImpl) CreateDialog(dialog *models.Dialog) (*models.Dialog, error) {
	if err := d.db.Create(&dialog).Error; err != nil {
		d.logger.Error("CreateDialog error", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrCreateDialog, err)
	}

	return dialog, nil
}

func NewDialogService(db *gorm.DB, logger *zap.Logger) DialogService {
	return &DialogServiceImpl{db, logger}
}
