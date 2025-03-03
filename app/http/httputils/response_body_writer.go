package httputils

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func NewResponseBodyWriter(responseWriter gin.ResponseWriter, body *bytes.Buffer) *ResponseBodyWriter {
	return &ResponseBodyWriter{
		ResponseWriter: responseWriter,
		body:           body,
	}
}

func (r ResponseBodyWriter) Body() *bytes.Buffer {
	return r.body
}
