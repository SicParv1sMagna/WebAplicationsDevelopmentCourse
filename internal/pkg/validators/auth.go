package validators

import (
	"errors"
	"project/internal/model"
)

func ValidateRegistrationData(user model.User) error {
	if user.FirstName == "" {
		return errors.New("имя должно быть заполнено")
	}

	if user.SecondName == "" {
		return errors.New("фамилия должна быть заполнена")
	}

	if user.Email == "" {
		return errors.New("поле почта должно быть заполнено")
	}

	if user.Password == "" {
		return errors.New("поле пароль должно быть заполнено")
	}

	if user.Password != user.RepeatPassword {
		return errors.New("пароли не совпадают")
	}

	if len(user.Password) > 20 {
		return errors.New("пароль не должен первышать 20 символов")
	}

	if len(user.Password) < 8 {
		return errors.New("пароль не должен быть меньше 8 символов")
	}

	return nil
}

func ValidateAuthorizationData(user model.User) error {
	if user.Password == "" {
		return errors.New("пароль не должен быть пустым")
	}

	if user.Email == "" {
		return errors.New("почта должна быть заполнена")
	}

	return nil
}
