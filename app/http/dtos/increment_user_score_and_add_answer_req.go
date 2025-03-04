package dtos

type IncreaseUserScoreAndAddAnswerRequest struct {
	ProblemId uint   `json:"problem" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}
