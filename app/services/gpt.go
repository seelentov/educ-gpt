package services

type GptService interface {
	AIService
}

type GptResponse struct {
	Choices []*GptResponseChoice `json:"choices"`
}

type GptResponseChoice struct {
	Message *GptResponseChoiceMessage `json:"message"`
}

type GptResponseChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
