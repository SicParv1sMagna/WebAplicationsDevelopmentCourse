package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (d *Delivery) GetContributor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении id"))
		return
	}

	contributor, markdowns, err := d.usecase.GetContributor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contributor": contributor,
		"markdown":    markdowns,
	})
}

func (d *Delivery) GetAllContributorsFromMarkdown(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	email := c.DefaultQuery("email", "")
	start_date := c.DefaultQuery("start_date", "")
	end_date := c.DefaultQuery("end_date", "")
	status := c.DefaultQuery("status", "")

	contributors, err := d.usecase.GetAllContributorsFromMarkdown(email, status, start_date, end_date, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, contributors)
}

func (d *Delivery) UpdateContributorAccessByModerator(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных").Error())
		return
	}

	if err := d.usecase.UpdateContributorAccessByModerator(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "статус изменен")
}

func (d *Delivery) UpdateContributroAccessByAdmin(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ошибка при получении данных"))
		return
	}

	if err := d.usecase.UpdateContributorAccessByAdmin(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "статус изменен")
}

func (d *Delivery) GetAllContirbutors(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	start_date := c.DefaultQuery("start_date", "")
	end_date := c.DefaultQuery("end_date", "")
	status := c.DefaultQuery("status", "")

	contributors, err := d.usecase.GetAllContributors(email, status, start_date, end_date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, contributors)
}

func (d *Delivery) UpdateContributorData(c *gin.Context) {
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
		err = d.usecase.Repository.UpdateContributorData(uint(id), email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при изменении данных"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "успешно"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "данные не введены"})
}
