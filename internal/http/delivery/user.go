package delivery

import (
	"errors"
	"net/http"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

// @Summary Зарегистрировать нового пользователя
// @Description Регистрация нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param userReq body model.UserRegisterReq true "Информация пользователя, необходимая для регистрации"
// @Success 201 {object} model.User "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func (d *Delivery) RegisterUser(c *gin.Context) {
	var userRegisterReq model.UserRegisterReq

	if err := c.BindJSON(&userRegisterReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	user, err := d.usecase.RegisterUser(userRegisterReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Авторизовать пользователя
// @Description Авторизация
// @Tags users
// @Accept json
// @Produce json
// @Param userReq body model.UserAuthReq true "Информация пользователя для авторизации"
// @Success 200 {object} string "Токен для авторизации"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (d *Delivery) LoginUser(c *gin.Context) {
	var userLoginReq model.UserAuthReq

	if err := c.BindJSON(&userLoginReq); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	token, err := d.usecase.AuthUser(userLoginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, token)
}

func (d *Delivery) GetMe(c *gin.Context) {
	userID := c.MustGet("UserID").(int)

	user, err := d.usecase.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// Обновление пользователем информации о себе
func UpdateUserInfo(c *gin.Context) {

}

// Удаление аккаунта пользователя
func DeleteUser(c *gin.Context) {

}
