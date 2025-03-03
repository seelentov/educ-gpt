package validator

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ParseValidationErrors(err error) map[string]string {
	var errors = make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		err.Error()
		f := strings.ToLower(e.Field())
		errors[f] = getErrorMsg(e.ActualTag(), e.Param())
	}
	return errors
}
