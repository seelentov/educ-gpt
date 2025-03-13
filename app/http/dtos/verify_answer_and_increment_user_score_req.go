package dtos

type VerifyAnswerAndIncrementUserScoreRequest struct {
	ProblemId uint   `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}
