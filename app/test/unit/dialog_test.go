package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"errors"
	"testing"
)

var (
	dialogSrv services.DialogService

	tempDialog  models.Dialog
	userID      uint = 1
	messageText      = "Hello world!"
)

func TestInitDialogService(t *testing.T) {
	dialogSrv = dic.DialogService()
}

func TestCanCreateDialog(t *testing.T) {
	dialog, err := dialogSrv.CreateDialog(&models.Dialog{
		UserID: userID,
	})
	if err != nil {
		t.Error(err)
		return
	}

	tempDialog = *dialog
}

func TestCanThrowMessage(t *testing.T) {
	if err := dialogSrv.AddDialogItem(&models.DialogItem{
		Text:     messageText,
		IsUser:   true,
		DialogID: tempDialog.ID,
	}); err != nil {
		t.Error(err)
		return
	}
}

func TestCanGetDialogs(t *testing.T) {
	dialogs, err := dialogSrv.GetDialogsByUserID(userID)
	if err != nil {
		t.Error(err)
		return
	}

	if len(dialogs) != 1 {
		t.Error("Expected 1 dialog got ", len(dialogs))
		return
	}
}

func TestCanGetDialog(t *testing.T) {
	dialog, err := dialogSrv.GetDialog(tempDialog.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if dialog == nil {
		t.Error("dialog is nil")
		return
	}

	if dialog.UserID != tempDialog.UserID {
		t.Errorf("Expected dialog.UserID to be %v but got %v", tempDialog.UserID, dialog.UserID)
		return
	}

	if len(dialog.DialogItems) != 1 {
		t.Error("Expected 1 message got ", len(dialog.DialogItems))
		return
	}

	if dialog.DialogItems[0].Text != messageText {
		t.Errorf("Expected message text to be %v but got %v", messageText, dialog.DialogItems[0].Text)
		return
	}
}

func TestCanRemoveDialog(t *testing.T) {
	err := dialogSrv.RemoveDialog(tempDialog.ID)
	if err != nil {
		t.Error(err)
	}

	dialog, err := dialogSrv.GetDialog(tempDialog.ID)
	if err != nil && !errors.Is(err, services.ErrDialogNotFound) {
		t.Error(err)
		return
	}

	if dialog != nil {
		t.Errorf("Expected dialog is nil but got %v", dialog)
		return
	}
}
