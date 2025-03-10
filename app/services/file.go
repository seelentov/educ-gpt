package services

import "mime/multipart"

type FileService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
	DeleteFile(path string) (bool, error)
}
