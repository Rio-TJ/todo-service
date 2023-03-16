package helpers

func ValidationMessageForTag(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	default:
		return ""
	}
}
