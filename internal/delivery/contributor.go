package delivery

import (
	"net/http"
	"project/internal/middleware"
	"project/internal/model"
	"project/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetContributor(repository *repository.Repository, c *gin.Context) {
	var contributor model.Contributor

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

	contributor, markdowns, err := repository.GetContributorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contributor": contributor,
		"markdown":    markdowns,
	})
}

func GetAllContributorsFromMarkdown(repository *repository.Repository, c *gin.Context) {
	var contributors []model.Contributor

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	email := c.DefaultQuery("email", "")
	start_date := c.DefaultQuery("start_date", "")
	end_date := c.DefaultQuery("end_date", "")
	status := c.DefaultQuery("status", "")

	contributors, err = repository.GetContributorsByMarkdownID(email, status, start_date, end_date, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, contributors)
}

func UpdateContributorAccess(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cid, _ := jsonData["Contributor_ID"].(float64)
	mid, _ := jsonData["Markdown_ID"].(float64)
	access, _ := jsonData["Access"].(string)

	err := repository.UpdateContributorAccess(uint(mid), uint(cid), access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, middleware.Response{
		Status:  "Success",
		Message: "Статус изменен",
	})
}

func GetAllContirbutors(repository *repository.Repository, c *gin.Context) {
	email := c.DefaultQuery("email", "")
	start_date := c.DefaultQuery("start_date", "")
	end_date := c.DefaultQuery("end_date", "")
	status := c.DefaultQuery("status", "")

	var contributors []model.Contributor

	contributors, err := repository.GetAllContributors(email, status, start_date, end_date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, contributors)
}

func UpdateContributorData(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id"})
		return
	}

	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if email, ok := jsonData["email"].(string); ok {
		err = repository.UpdateContributorData(uint(id), email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при изменении данных"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "успешно"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "данные не введены"})
}
