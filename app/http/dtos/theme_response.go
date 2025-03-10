package dtos

import "educ-gpt/models"

type ThemeResponse struct {
	Text     string            `json:"text"`
	Problems []*models.Problem `json:"problems,omitempty"`
}
