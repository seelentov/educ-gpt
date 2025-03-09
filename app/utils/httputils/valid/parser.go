package valid

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func ParseValidationErrors(err error) map[string]string {
	var errorsMap = make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		err.Error()
		f := strcase.ToSnake(e.Field())
		errorsMap[f] = getErrorMsg(e.ActualTag(), e.Param())
	}
	return errorsMap
}
