package validators

import (
	"project/internal/middleware"
	"project/internal/model"
)

func ValidateMarkdown(md model.Markdown) middleware.Response {
	if md.Name == "" {
		return middleware.Response{
			Status:  "Failed",
			Message: "Название должно быть заполнено",
		}
	}

	if len(md.Name) > 26 {
		return middleware.Response{
			Status:  "Failed",
			Message: "Название должно быть не более 25 символов",
		}
	}

	return middleware.Response{
		Status:  "Success",
		Message: "",
	}
}
