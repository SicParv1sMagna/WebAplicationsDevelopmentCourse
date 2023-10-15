package delivery

import (
	"fmt"
	"net/http"
	"project/internal/middleware"
	jwttoken "project/internal/middleware/jwt"
	"project/internal/middleware/validators"
	"project/internal/model"
	"project/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMarkdown(repository *repository.Repository, c *gin.Context) {
	var markdown model.Markdown

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.Response{
			Status:  "Failed",
			Message: "Unauthorized",
		})
		return
	}

	// Достаем данные из JSON'а из запроса
	if err := c.BindJSON(&markdown); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(markdown)

	//	Валиидруем название markdown'а
	if err := validators.ValidateMarkdown(markdown); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	markdown.User_ID = uint(userID)

	if err := repository.CreateMarkdown(markdown); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, middleware.Response{
		Status:  "Created",
		Message: "MD создан",
	})
}

func GetMarkdown(repository *repository.Repository, c *gin.Context) {
	var md model.Markdown

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, middleware.Response{
			Status:  "Failed",
			Message: "id не может быть отрицательным",
		})
		return
	}

	md, err = repository.GetMarkdownById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, md)
}

func GetAllMarkdowns(repository *repository.Repository, c *gin.Context) {
	var md []model.Markdown

	md, err := repository.GetAllMarkdowns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, md)
}

func DeleteMarkdown(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, middleware.Response{
			Status:  "Failed",
			Message: "id не может быть отрицательным",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = repository.DeleteMarkdownById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
}

func UpdateMarkdown(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mdID, idOk := jsonData["Markdown_ID"].(float64)
	Name, nameOk := jsonData["Name"].(string)
	Content, contentOk := jsonData["Content"].(string)

	if !idOk || mdID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing Markdown_ID"})
		return
	}

	candidate, err := repository.GetMarkdownById(uint(mdID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if nameOk {
		candidate.Name = Name
	}

	if contentOk {
		candidate.Content = Content
	}

	err = repository.UpdateMarkdownById(candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Markdown updated successfully",
	})
}

func RequestContribution(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, middleware.Response{
			Status:  "Failed",
			Message: "id не может быть отрицательным",
		})
	}

	var contributor model.Contributor
	if err := c.BindJSON(&contributor); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = repository.RequestContribution(contributor, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, middleware.Response{
		Status:  "Success",
		Message: "заявка отправлена",
	})
}
