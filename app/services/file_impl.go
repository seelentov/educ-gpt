package services

import (
	"educ-gpt/utils/securityutils"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	ErrCompressFile = errors.New("compress file error")
	ErrSaveFile     = errors.New("save file error")
	ErrDeleteFile   = errors.New("delete file error")
)

type FileServiceImpl struct {
	logger        *zap.Logger
	storagePrefix string
}

func (f FileServiceImpl) UploadImage(file *multipart.FileHeader) (string, error) {
	if file.Size > maxImageFileSize {
		compressedFile, err := f.CompressImage(file)
		if err != nil {
			f.logger.Error("compress file error", zap.Error(err))
			return "", fmt.Errorf("%w:%w", ErrCompressFile, err)
		}
		file = compressedFile
	}

	fileName := securityutils.GenerateKey(50)
	fileExt := filepath.Ext(file.Filename)
	newFileName := fileName + fileExt

	filePath := filepath.Join("app", "storage", newFileName)

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

func (f FileServiceImpl) CompressImage(file *multipart.FileHeader) (*multipart.FileHeader, error) {
	src, err := file.Open()
	if err != nil {
		f.logger.Error("open file error", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
	}
	defer src.Close()

	img, format, err := image.Decode(src)
	if err != nil {
		f.logger.Error("decode image error", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
	}

	quality := 100
	for {
		tempFile, err := os.CreateTemp("", "compressed-*.jpg")
		if err != nil {
			f.logger.Error("create temp file error", zap.Error(err))
			return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
		}
		defer tempFile.Close()

		switch format {
		case "jpeg", "jpg":
			if err := jpeg.Encode(tempFile, img, &jpeg.Options{Quality: quality}); err != nil {
				f.logger.Error("encode image error", zap.Error(err))
				return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
			}
		case "png":
			if err := png.Encode(tempFile, img); err != nil {
				f.logger.Error("encode image error", zap.Error(err))
				return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
			}
		default:
			f.logger.Error("unknown image format", zap.String("format", format))
			return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
		}

		fileInfo, err := tempFile.Stat()
		if err != nil {
			f.logger.Error("get file info error", zap.Error(err))
			return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
		}

		if fileInfo.Size() <= maxImageFileSize {
			return &multipart.FileHeader{
				Filename: file.Filename,
				Size:     fileInfo.Size(),
				Header:   file.Header,
			}, nil
		}

		quality -= 10
		if quality <= 0 {
			f.logger.Error("failed to compress to required size", zap.Int("quality", quality))
			return nil, fmt.Errorf("%w:%w", ErrCompressFile, err)
		}

		img = resize.Resize(uint(img.Bounds().Dx()/2), uint(img.Bounds().Dy()/2), img, resize.Lanczos3)
	}
}

func NewFileServiceImpl(logger *zap.Logger, storagePrefix string) FileService {
	return &FileServiceImpl{logger, storagePrefix}
}
