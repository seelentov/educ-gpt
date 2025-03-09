package dtos

type LoginRequest struct {
	Credential string `json:"credential" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
