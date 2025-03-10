package dtos

type ValidationErrorResponse struct {
	Error map[string]string `json:"error"`
}
