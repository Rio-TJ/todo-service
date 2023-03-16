package helpers

func ValidationMessageForTag(tag, param string) string {
	switch tag {
	case "required":
		return "Это поле обязательно к заполнению"
	default:
		return ""
	}
}
