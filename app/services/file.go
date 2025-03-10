package services

import "mime/multipart"

const maxImageFileSize = 1048576

type FileService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
	DeleteFile(path string) (bool, error)
	CompressImage(file *multipart.FileHeader) (*multipart.FileHeader, error)
}
