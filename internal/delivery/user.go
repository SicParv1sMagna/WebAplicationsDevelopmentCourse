package delivery

import (
	"net/http"
	"project/internal/middleware"
	"project/internal/middleware/authorization"
	jwttoken "project/internal/middleware/jwt"
	"project/internal/middleware/validators"
	"project/internal/model"
	"project/internal/repository"

	"github.com/gin-gonic/gin"
)

// Создание нового пользователя
func RegisterUser(repository *repository.Repository, c *gin.Context) {
	var user model.User

	// Достаем данные из JSON'а из запроса
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Валидация введенной пользователем информации
	if err := validators.ValidateRegistrationData(user); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Пробуем достать из базы пользователя с таким же Email'ом
	candidate, err := repository.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Проверяем, зарегестрирован ли пользователь
	if candidate == user {
		c.JSON(http.StatusInternalServerError, middleware.Response{
			Status:  "Failed",
			Message: "Такой пользователь уже существует",
		})
		return
	}

	// Хэшируем пароль с использованием bcrypt
	user.Password, err = authorization.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Создаем пользователя
	err = repository.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, middleware.Response{
		Status:  "Created",
		Message: "Пользователь зарегестрирован",
	})
}

// Авторизация пользователя
func AuthUser(repository *repository.Repository, c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := validators.ValidateAuthorizationData(user); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	candidate, err := repository.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if ok := authorization.CheckPasswordHash(user.Password, candidate.Password); !ok {
		c.JSON(http.StatusBadRequest, middleware.Response{
			Status:  "Failed",
			Message: "Пароли не совпадают",
		})
		return
	}

	tokenString, err := jwttoken.GenerateJWTToken(uint(candidate.User_ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Value:    tokenString,
		HttpOnly: true,
	})

	// c.JSON(http.StatusOK, middleware.Response{
	// 	Status:  "Success",
	// 	Message: "Авторизован",
	// })

	c.JSON(http.StatusOK, tokenString)
}

// Обновление пользователем информации о себе
func UpdateUserInfo(c *gin.Context) {

}

// Удаление аккаунта пользователя
func DeleteUser(c *gin.Context) {

}
