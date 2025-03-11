package daemons

import (
	"context"
	"educ-gpt/models"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ClearUnusedFilesDaemon struct {
	db     *gorm.DB
	logger *zap.Logger
	sleep  time.Duration
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *ClearUnusedFilesDaemon) Work() {
	var users []*models.User
	c.db.Select("avatar_url").Find(&users)

	filesIsUse := make([]string, 0)
	for i := range users {
		filesIsUse = append(filesIsUse, users[i].AvatarUrl)
	}

	storageDir := "storage"
	filesToDelete, err := c.getUnusedFiles(storageDir, filesIsUse)
	if err != nil {
		c.logger.Error("error getting unused files", zap.Error(err))
		return
	}

	for _, file := range filesToDelete {
		err := os.Remove(file)
		if err != nil {
			c.logger.Error("error removing file", zap.Error(err))
			continue
		}
	}
}

func (c *ClearUnusedFilesDaemon) getUnusedFiles(storageDir string, filesIsUse []string) ([]string, error) {
	unusedFiles := make([]string, 0)

	err := filepath.Walk(storageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			c.logger.Error("error walking", zap.Error(err))
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !c.isFileUsed(path, filesIsUse) {
			unusedFiles = append(unusedFiles, path)
		}

		return nil
	})

	if err != nil {
		c.logger.Error("error walking", zap.Error(err))
		return nil, err
	}

	return unusedFiles, nil
}

func (c *ClearUnusedFilesDaemon) isFileUsed(filePath string, filesIsUse []string) bool {
	fileName := filepath.Base(filePath)

	for _, usedFile := range filesIsUse {
		if strings.Contains(usedFile, fileName) {
			return true
		}
	}

	return false
}

func (c *ClearUnusedFilesDaemon) Start() {
	go func() {
		c.logger.Info("ClearUnusedFilesDaemon started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearUnusedFilesDaemon works")
				c.Work()
				c.logger.Info("ClearUnusedFilesDaemon works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearUnusedFilesDaemon stopped")
				return
			}
		}
	}()
}

func (c *ClearUnusedFilesDaemon) Stop() {
	c.logger.Info("Stopping ClearUnusedFilesDaemon...")
	c.cancel()
}

func NewClearUnusedFilesDaemon(
	db *gorm.DB,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearUnusedFilesDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearUnusedFilesDaemon{
			db:     db,
			logger: logger,
			sleep:  sleep,
			ctx:    ctx,
			cancel: cancel,
		},
		ctx,
		cancel
}
