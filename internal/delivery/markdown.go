package delivery

import (
	"fmt"
	"io"
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
		c.JSON(http.StatusUnauthorized, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//	Валиидруем название markdown'а
	if err := validators.ValidateMarkdown(markdown); err.Status == "Failed" {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	markdown.User_ID = uint(userID)

	if err := repository.CreateMarkdown(markdown); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, md)
}

func GetAllMarkdowns(repository *repository.Repository, c *gin.Context) {
	var md []model.Markdown

	name := c.DefaultQuery("name", "")

	md, err := repository.GetAllMarkdowns(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, md)
}

func DeleteMarkdown(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = repository.DeleteMarkdownById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
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

func SearchMarkdown(repository *repository.Repository, c *gin.Context) {
	query := c.DefaultQuery("query", "")

	md, err := repository.SearchMarkdown(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, md)
}

func AddMarkdownToContributor(repository *repository.Repository, c *gin.Context) {
	markdownID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, middleware.Response{
			Status:  "Failed",
			Message: "Invalid ID",
		})
	}

	contributor, markdowns, err := repository.AddMdToLastReader(uint(markdownID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, middleware.Response{
			Status:  "Failed",
			Message: "Error occured",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"contributor": contributor,
		"markdowns":   markdowns,
	})
}

func DeleteContributorFromMd(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cid, cidOk := jsonData["Contributor_ID"].(float64)
	mid, midOk := jsonData["Markdown_ID"].(float64)

	if !cidOk || !midOk || mid <= 0 || cid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
	}

	fmt.Println(cid, mid)

	err := repository.DeleteContributorFromMd(uint(mid), uint(cid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, middleware.Response{
		Status:  "Success",
		Message: "Success",
	})
}

func AddMarkdownIcon(repository *repository.Repository, c *gin.Context) {
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

	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение"})
		return
	}

	contentType := image.Header.Get("Content-Type")

	err = repository.AddMarkdownIcon(id, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "изображение успешно загружено"})
}
