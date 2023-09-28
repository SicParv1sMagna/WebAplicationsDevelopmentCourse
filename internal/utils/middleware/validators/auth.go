package validators

import (
	"project/internal/model"
	"project/internal/utils"
	"unicode"
)

func ValidateRegistrationData(user model.User) *utils.Response {
	if user.FirstName == "" {
		return &utils.Response{
			Status:  "Failed",
			Message: "Поле имя должно быть заполнено",
		}
	}

	if user.SecondName == "" {
		return &utils.Response{
			Status:  "Failed",
			Message: "Поле фамилия должно быть заполнено",
		}
	}

	if user.Email == "" {
		return &utils.Response{
			Status:  "Failed",
			Message: "Поле Email должно быть заполнено",
		}
	}

	if user.Password == "" {
		return &utils.Response{
			Status:  "Failed",
			Message: "Поле пароль должно быть заполнено",
		}
	}

	if user.Password != user.RepeatPassword {
		return &utils.Response{
			Status:  "Failed",
			Message: "Пароли не совпадают",
		}
	}

	if len(user.Password) > 20 {
		return &utils.Response{
			Status:  "Failed",
			Message: "Пароль должен быть меньше 20 символов",
		}
	}

	if len(user.Password) < 8 {
		return &utils.Response{
			Status:  "Failed",
			Message: "Пароль должен быть больше 8 символов",
		}
	}

	hasLetter := false
	hasDigit := false

	for _, char := range user.Password {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasLetter = true
		}
		if hasDigit && hasLetter {
			return &utils.Response{
				Status:  "Failed",
				Message: "Пароль должен содержать цифры и буквы",
			}
		}
	}

	return &utils.Response{
		Status:  "Success",
		Message: "",
	}
}
