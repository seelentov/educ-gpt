package services

import (
	"educ-gpt/models"
	"errors"
)

var (
	ErrDialogNotFound = errors.New("dialog not found")
	ErrGetDialogs     = errors.New("get dialogs error")
	ErrAddDialogItem  = errors.New("add dialog item error")
	ErrCreateDialog   = errors.New("create dialog error")
	ErrDeleteDialog   = errors.New("create delete error")
)

type DialogService interface {
	GetDialog(dialogID uint) (*models.Dialog, error)
	AddDialogItem(dialogItem *models.DialogItem) error
	RemoveDialog(dialogID uint) error
	CreateDialog(dialog *models.Dialog) (*models.Dialog, error)
	GetDialogsByUserID(userID uint) ([]*models.Dialog, error)
}
