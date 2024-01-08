package delivery

import (
	"errors"
	"net/http"
	"project/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (d *Delivery) CreateMarkdown(c *gin.Context) {
	var markdownReq model.Markdown

	userID := c.MustGet("UserID").(int)

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
	userID := c.MustGet("userID").(int)
	name := c.DefaultQuery("name", "")

	md, id, err := d.usecase.GetAllMarkdown(name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"Contributor_id": id,
		"Markdowns": md})
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

	userID := c.MustGet("UserID").(int)

	if err = d.usecase.AddMarkdownToContributor(uint(markdownID), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "добавлен")
}

func (d *Delivery) DeleteContributorFromMd(c *gin.Context) {
	userID := c.MustGet("UserID").(int)

	markdownID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении данных"})
		return
	}

	if err := d.usecase.DeleteContributorFromMd(markdownID, userID); err != nil {
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

	url, err := d.usecase.AddMarkdownIcon(uint(id), image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, url)
}

func (d *Delivery) RequestContribution(c *gin.Context) {
	userID := c.MustGet("UserID").(int)

	markdownID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	if err = d.usecase.RequestContribution(uint(userID), uint(markdownID)); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "заявка успешно создана")
}
