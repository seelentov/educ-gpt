package dtos

type IncreaseUserScoreAndAddAnswerRequest struct {
	ProblemId uint   `json:"problem_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}
