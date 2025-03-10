package dtos

type ErrorResponse struct {
	Error string `json:"error"`
}

func UnauthorizedResponse() *ErrorResponse {
	return &ErrorResponse{Error: "Unauthorized"}
}

func NotFoundResponse() *ErrorResponse {
	return &ErrorResponse{Error: "Not found"}
}

func InternalServerErrorResponse() *ErrorResponse {
	return &ErrorResponse{Error: "Internal server error"}
}
