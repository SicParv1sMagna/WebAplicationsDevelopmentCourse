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

	//	Валиидруем название markdown'а
	if err := validators.ValidateMarkdown(markdown); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	markdown.Moderator_ID = uint(userID)

	if err := repository.CreateMarkdown(markdown, userID); err != nil {
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

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.Response{
			Status:  "Failed",
			Message: "Unauthorized",
		})
		return
	}

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

	md, err = repository.GetMarkdownById(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, md)
}

func GetAllMarkdowns(repository *repository.Repository, c *gin.Context) {
	var md []model.Markdown

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.Response{
			Status:  "Failed",
			Message: "Unauthorized",
		})
		return
	}

	md, err = repository.GetAllMarkdowns(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, md)
}

func DeleteMarkdown(repository *repository.Repository, c *gin.Context) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.Response{
			Status:  "Failed",
			Message: "Unauthorized",
		})
		return
	}

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

	err = repository.DeleteMarkdownById(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
}

func UpdateMarkdown(repository *repository.Repository, c *gin.Context) {

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	userID, err := jwttoken.GetUserIDbyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.Response{
			Status:  "Failed",
			Message: "Unauthorized",
		})
		return
	}

	mdID, err := strconv.Atoi(c.Query("Markdown_ID"))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	Name := c.Query("Name")
	Content := c.Query("Name")

	candidate, err := repository.GetMarkdownById(uint(mdID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if Name != "" {
		candidate.Name = Name
	}

	if Content != "" {
		candidate.Content = Content
	}

	err = repository.UpdateMarkdownById(candidate, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
}
