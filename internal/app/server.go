package app

import (
	"log"
	"project/internal/delivery"

	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server starting")

	// Создаем роутинг
	router := gin.Default()

	//	http://localhost:8080/user
	user := router.Group("/user")
	{
		//	http://localhost:8080/user/register
		user.POST("/register", func(c *gin.Context) {
			delivery.RegisterUser(a.repository, c)
		})

		//	http://localhost:8080/user/login
		user.POST("/login", func(c *gin.Context) {
			delivery.AuthUser(a.repository, c)
		})

		//	http://localhost:8080/user/delete
		user.DELETE("/delete", delivery.DeleteUser)

		//	http://localhost:8080/user/edit-info
		user.PUT("/edit-info", delivery.UpdateUserInfo)
	}

	// TODO: Add contributor routes

	// contributor := router.Group("/contribute")
	// {

	// }

	notes := router.Group("/notes")
	{
		markdowns := notes.Group("/markdown")
		{
			markdowns.POST("/create-markdown", delivery.CreateMarkdown)

			markdowns.GET("/get-markdown", delivery.GetMarkdown)

			markdowns.GET("/get-all-markdowns", delivery.GetAllMarkdowns)

			markdowns.DELETE("/delete-markdown", delivery.DeleteMarkdown)

			markdowns.PUT("/update-markdown", delivery.UpdateMarkdown)
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
