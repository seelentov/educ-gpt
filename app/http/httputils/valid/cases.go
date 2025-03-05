package valid

import "fmt"

func getErrorMsg(tag string, param string) string {
	switch tag {
	case "required":
		return "Is required"
	case "email":
		return "Should be valid email"
	case "lte":
		return "Should be less than " + param
	case "gte":
		return "Should be greater than " + param
	case "number":
		return "Should be valid phone number"
	}

	return fmt.Sprintf("%s:%s", tag, param)
}
