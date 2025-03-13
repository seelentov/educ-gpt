package dtos

type VerifyAnswerRequest struct {
	Answer   string `json:"answer" binding:"required"`
	Problem  string `json:"problem" binding:"required"`
	Language string `json:"language" binding:"required"`
}
