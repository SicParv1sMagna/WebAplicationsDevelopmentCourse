package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"project/internal/model"
	jwttoken "project/internal/pkg/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (d *Delivery) CreateMarkdown(c *gin.Context) {
	var markdownReq model.Markdown

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("перед созданием необходимо авторизоваться"))
		return
	}

	// Достаем данные из JSON'а из запроса
	if err := c.BindJSON(&markdownReq); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	markdown, err := d.usecase.CreateMarkdown(markdownReq, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, markdown)
}

func (d *Delivery) GetAllMarkdowns(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	md, err := d.usecase.GetAllMarkdown(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, md)
}

func (d *Delivery) GetMarkdown(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	markdown, err := d.usecase.GetMarkdown(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, markdown)
}

func (d *Delivery) DeleteMarkdown(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = d.usecase.DeleteMarkdown(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "удален")
}

func (d *Delivery) UpdateMarkdown(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	if err := d.usecase.UpdateMarkdown(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "маркдаун обновлен")
}

func (d *Delivery) AddMarkdownToContributor(c *gin.Context) {
	markdownID, err := strconv.Atoi(c.Param("markdown_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	contributorID, err := strconv.Atoi(c.Param("contributor_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	fmt.Println(markdownID, contributorID)
	if err = d.usecase.AddMarkdownToContributor(uint(markdownID), uint(contributorID)); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "добавлен")
}

func (d *Delivery) DeleteContributorFromMd(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := d.usecase.DeleteContributorFromMd(jsonData); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	c.JSON(http.StatusOK, "удален")
}

func (d *Delivery) AddMarkdownIcon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоступный ID багажа"})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
		return
	}

	if err = d.usecase.AddMarkdownIcon(uint(id), image); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "иконка добавлена")
}

func (d *Delivery) RequestContribution(c *gin.Context) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("неавторизован"))
		return
	}

	markdownID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	if err = d.usecase.RequestContribution(userID, uint(markdownID)); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "заявка успешно создана")
}
