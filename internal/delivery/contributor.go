package delivery

import (
	"fmt"
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

	contributors, err = repository.GetContributorsByMarkdownID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, contributors)
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

func UpdateContributorAccess(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cid, _ := jsonData["Contributor_ID"].(float64)
	mid, _ := jsonData["Markdown_ID"].(float64)
	access, _ := jsonData["Access"].(string)

	fmt.Println(cid, mid, access)

	err := repository.UpdateContributorAccess(uint(mid), uint(cid), access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, middleware.Response{
		Status:  "Success",
		Message: "Success",
	})
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
