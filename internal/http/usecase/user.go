package usecase

import (
	"errors"
	"project/internal/model"
	"project/internal/pkg/authorization"
	jwttoken "project/internal/pkg/jwt"
	"project/internal/pkg/roles"
	"project/internal/pkg/validators"
)

func (uc *UseCase) RegisterUser(userRegReq model.UserRegisterReq) (model.User, error) {
	if err := validators.ValidateRegistrationData(userRegReq); err != nil {
		return model.User{}, err
	}

	candidate, err := uc.Repository.GetUserByEmail(userRegReq.Email)
	if err != nil {
		return model.User{}, err
	}

	if candidate.Email == userRegReq.Email {
		return model.User{}, errors.New("такой пользователь уже существует")
	}

	userRegReq.Password, err = authorization.HashPassword(userRegReq.Password)
	if err != nil {
		return model.User{}, errors.New("ошибка шифрования пароля")
	}

	var user model.User = model.User{
		FirstName:      userRegReq.FirstName,
		SecondName:     userRegReq.SecondName,
		MiddleName:     userRegReq.MiddleName,
		Email:          userRegReq.Email,
		Password:       userRegReq.Password,
		RepeatPassword: userRegReq.RepeatPassword,
		Role:           roles.User,
	}

	err = uc.Repository.CreateUser(user)
	if err != nil {
		return model.User{}, errors.New("ошибка создания пользователя")
	}

	return user, nil
}

func (uc *UseCase) AuthUser(userLoginReq model.UserAuthReq) (string, error) {
	if err := validators.ValidateAuthorizationData(userLoginReq); err != nil {
		return "", err
	}

	candidate, err := uc.Repository.GetUserByEmail(userLoginReq.Email)
	if err != nil {
		return "", errors.New("такого пользователя нету")
	}

	if ok := authorization.CheckPasswordHash(userLoginReq.Password, candidate.Password); !ok {
		return "", errors.New("пароли не совпадают")
	}

	tokenString, err := jwttoken.GenerateJWTToken(uint(candidate.User_ID), int(candidate.Role))
	if err != nil {
		return "", errors.New("ошибка при авторизации пользователя")
	}

	if err = uc.Repository.SaveJWTToken(uint(candidate.User_ID), tokenString); err != nil {
		return "", errors.New("ошибка при сохранении токена авторизации")
	}

	return tokenString, nil
}

func (uc *UseCase) GetMe(userID int) (model.User, error) {
	if userID <= 0 {
		return model.User{}, errors.New("id не может быть отрицательным")
	}

	user, err := uc.Repository.GetUserById(uint(userID))
	if err != nil {
		return model.User{}, errors.New("ошибка при получении данных")
	}

	return user, nil
}
