package usecase

import (
	"errors"
	"project/internal/model"
	"project/internal/pkg/authorization"
	jwttoken "project/internal/pkg/jwt"
	"project/internal/pkg/validators"
)

func (uc *UseCase) RegisterUser(user model.User) (model.User, error) {
	if err := validators.ValidateRegistrationData(user); err != nil {
		return model.User{}, err
	}

	candidate, err := uc.Repository.GetUserByEmail(user.Email)
	if err != nil {
		return model.User{}, err
	}

	if candidate == user {
		return model.User{}, errors.New("такой пользователь уже существует")
	}

	user.Password, err = authorization.HashPassword(user.Password)
	if err != nil {
		return model.User{}, errors.New("ошибка шифрования пароля")
	}

	err = uc.Repository.CreateUser(user)
	if err != nil {
		return model.User{}, errors.New("ошибка создания пользователя")
	}

	return user, nil
}

func (uc *UseCase) AuthUser(user model.User) (string, error) {
	if err := validators.ValidateAuthorizationData(user); err != nil {
		return "", err
	}

	candidate, err := uc.Repository.GetUserByEmail(user.Email)
	if err != nil {
		return "", errors.New("такого пользователя нету")
	}

	if ok := authorization.CheckPasswordHash(user.Password, candidate.Password); !ok {
		return "", errors.New("пароли не совпадают")
	}

	tokenString, err := jwttoken.GenerateJWTToken(uint(candidate.User_ID))
	if err != nil {
		return "", errors.New("ошибка при авторизации пользователя")
	}

	if err = uc.Repository.SaveJWTToken(uint(candidate.User_ID), tokenString); err != nil {
		return "", errors.New("ошибка при сохранении токена авторизации")
	}

	return tokenString, nil
}
