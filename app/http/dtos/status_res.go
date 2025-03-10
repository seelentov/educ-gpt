package dtos

type StatusResponse struct {
	Status string `json:"status"`
}

func OkResponse() *StatusResponse {
	return &StatusResponse{Status: "ok"}
}
