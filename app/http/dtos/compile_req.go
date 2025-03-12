package dtos

type CompileRequest struct {
	Code string `json:"code" binding:"required"`
}
