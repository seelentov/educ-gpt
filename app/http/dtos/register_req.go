package dtos

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,lte=100"`
	Email    string `json:"email" binding:"required,email,lte=100"`
	Number   string `json:"number" binding:"required,gte=8,lte=100,number"`
	Password string `json:"password" binding:"required,gte=8"`
}
