package impl

import (
	"educ-gpt/services"
	"educ-gpt/utils/securityutils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	ErrCompressFile       = errors.New("compress file error")
	ErrCompressFileEnough = errors.New("compress file enough error")
	ErrSaveFile           = errors.New("save file error")
	ErrDeleteFile         = errors.New("delete file error")
)

type FileServiceImpl struct {
	logger        *zap.Logger
	storagePrefix string
}

func (f FileServiceImpl) UploadImage(file *multipart.FileHeader) (string, error) {
	var newFileName string
	var filePath string

	for {
		fileName := securityutils.GenerateKey(50)
		fileExt := filepath.Ext(file.Filename)
		newFileName = fileName + fileExt
		filePath = filepath.Join("storage", newFileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
	}

	srcFile, err := file.Open()
	if err != nil {
		f.logger.Error("open file error", zap.Error(err))
		return "", fmt.Errorf("%w:%w", ErrSaveFile, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePath)
	if err != nil {
		f.logger.Error("create file error", zap.Error(err))
		return "", fmt.Errorf("%w:%w", ErrSaveFile, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		f.logger.Error("copy file error", zap.Error(err))
		return "", fmt.Errorf("%w:%w", ErrSaveFile, err)
	}

	return f.storagePrefix + "/" + newFileName, nil
}

func (f FileServiceImpl) DeleteFile(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f.logger.Error("file not exist", zap.String("path", filePath))
		return false, nil
	} else if err != nil {
		f.logger.Error("get file error", zap.Error(err))
		return false, fmt.Errorf("%w:%w", ErrDeleteFile, err)
	}

	if err := os.Remove(filePath); err != nil {
		f.logger.Error("delete file error", zap.Error(err))
		return false, fmt.Errorf("%w:%w", ErrDeleteFile, err)
	}

	return true, nil
}

func NewFileServiceImpl(logger *zap.Logger, storagePrefix string) services.FileService {
	return &FileServiceImpl{logger, storagePrefix}
}
