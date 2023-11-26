package delivery

import (
	"errors"
	"net/http"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

// Создание нового пользователя
func (d *Delivery) RegisterUser(c *gin.Context) {
	var userReq model.User

	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	user, err := d.usecase.RegisterUser(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Авторизация пользователя
func (d *Delivery) LoginUser(c *gin.Context) {
	var userReq model.User

	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	token, err := d.usecase.AuthUser(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Value:    token,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, token)
}

// Обновление пользователем информации о себе
func UpdateUserInfo(c *gin.Context) {

}

// Удаление аккаунта пользователя
func DeleteUser(c *gin.Context) {

}
