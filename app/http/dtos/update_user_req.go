package dtos

import "mime/multipart"

type UpdateUserRequest struct {
	Name       string                `form:"name,omitempty" json:"name,omitempty" binding:"omitempty,lte=100"`
	AvatarFile *multipart.FileHeader `form:"avatar_file,omitempty" json:"avatar_file,omitempty" binding:"omitempty"`
}
