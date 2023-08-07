package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server is starting")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main Website",
		})
	})

	r.GET("/wishes", func(c *gin.Context) {
		wishes := []string{
			"Хочу AKKO V3 Cream Blue switches",
			"Хочу Монитор 54 дюйма",
			"Хочу Porsche 991 GTR",
		}

		c.HTML(http.StatusOK, "wishes.tmpl", gin.H{
			"Wishes": wishes,
		})
	})

	r.Static("/image", "./resources")

	r.Run()

	log.Println("Server is shutting down")
}
