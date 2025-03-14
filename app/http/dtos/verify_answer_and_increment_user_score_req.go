package dtos

type VerifyAnswerAndIncrementUserScoreRequest struct {
	ProblemId uint   `json:"problem_id" binding:"required"`
	Language  string `json:"language"`
	Answer    string `json:"answer" binding:"required"`
}
