package valid

import "fmt"

func getErrorMsg(tag string, param string) string {
	switch tag {
	case "required":
		return "Обязателен"
	case "email":
		return "Должен быть валидным email"
	case "lte":
		return "Должен быть меньше чем " + param
	case "gte":
		return "Должен быть больше чем " + param
	case "number":
		return "Должен быть валидным номером телефона"
	}

	return fmt.Sprintf("%s:%s", tag, param)
}
