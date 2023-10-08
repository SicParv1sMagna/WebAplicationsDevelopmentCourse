package validators

import (
	"project/internal/middleware"
	"project/internal/model"
)

func ValidateRegistrationData(user model.User) *middleware.Response {
	if user.FirstName == "" {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Поле имя должно быть заполнено",
		}
	}

	if user.SecondName == "" {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Поле фамилия должно быть заполнено",
		}
	}

	if user.Email == "" {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Поле Email должно быть заполнено",
		}
	}

	if user.Password == "" {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Поле пароль должно быть заполнено",
		}
	}

	if user.Password != user.RepeatPassword {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Пароли не совпадают",
		}
	}

	if len(user.Password) > 20 {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Пароль должен быть меньше 20 символов",
		}
	}

	if len(user.Password) < 8 {
		return &middleware.Response{
			Status:  "Failed",
			Message: "Пароль должен быть больше 8 символов",
		}
	}

	// hasLetter := false
	// hasDigit := false

	// for _, char := range user.Password {
	// 	if unicode.IsDigit(char) {
	// 		hasDigit = true
	// 	} else {
	// 		hasLetter = true
	// 	}
	// 	if hasDigit && hasLetter {
	// 		return &middleware.Response{
	// 			Status:  "Failed",
	// 			Message: "Пароль должен содержать цифры и буквы",
	// 		}
	// 	}
	// }

	return &middleware.Response{
		Status:  "Success",
		Message: "",
	}
}

func ValidateAuthorizationData(user model.User) middleware.Response {
	if user.Password == "" {
		return middleware.Response{
			Status:  "Failed",
			Message: "Введите пароль",
		}
	}

	if user.Email == "" {
		return middleware.Response{
			Status:  "Failed",
			Message: "Введите Email",
		}
	}

	return middleware.Response{
		Status:  "Success",
		Message: "",
	}
}
