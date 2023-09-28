package delivery

import (
	"log"
	"net/http"
	"project/internal/model"
	"project/internal/repository"
	"project/internal/utils"
	"project/internal/utils/middleware/authorization"
	"project/internal/utils/middleware/validators"

	"github.com/gin-gonic/gin"
)

// Создание нового пользователя
func RegisterUser(repository *repository.Repository, c *gin.Context) {
	var user model.User

	// Достаем данные из JSON'а из запроса
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
	}

	// Валидация введенной пользователем информации
	if err := validators.ValidateRegistrationData(user); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
	}

	// Пробуем достать из базы пользователя с таким же Email'ом
	candidate, err := repository.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	// Проверяем, зарегестрирован ли пользователь
	if candidate == user {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Failed",
			Message: "Такой пользователь уже существует",
		})
	}

	// Хэшируем пароль с использованием bcrypt
	user.Password, err = authorization.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	// Создаем пользователя
	err = repository.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, utils.Response{
		Status:  "Created",
		Message: "Пользователь зарегестрирован",
	})
}

// Авторизация пользователя
func AuthUser(repository *repository.Repository, c *gin.Context) {

}

// Обновление пользователем информации о себе
func UpdateUserInfo(c *gin.Context) {}

// Удаление аккаунта пользователя
func DeleteUser(c *gin.Context) {}
