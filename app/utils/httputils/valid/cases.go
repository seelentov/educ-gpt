package valid

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
	case "url":
		return "Должен быть URL адресом"
	}

	result := tag

	if param != "" {
		result += ":" + param
	}

	return result
}
