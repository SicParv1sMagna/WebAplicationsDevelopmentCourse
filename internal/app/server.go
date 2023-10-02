package app

import (
	"log"
	"project/internal/delivery"
	jwttoken "project/internal/middleware/jwt"

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

	}

	//	http://localhost:8080/api
	api := router.Group("/api")
	{
		api.Use(jwttoken.CheckJWTToken())
		//	http://localhost:8080/api/user
		user := api.Group("/user")
		{
			//	http://localhost:8080/api/user/delete
			user.DELETE("/delete", delivery.DeleteUser)

			//	http://localhost:8080/api/user/edit-info
			user.PUT("/edit-info", delivery.UpdateUserInfo)
		}

		//	http://localhost:8080/api/notes
		notes := api.Group("/notes")
		{
			//	http://localhost:8080/api/notes/markdown
			markdowns := notes.Group("/markdown")
			{
				//	http://localhost:8080/api/notes/markdown/create-markdown
				markdowns.POST("/create-markdown", func(c *gin.Context) {
					delivery.CreateMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/get-markdown
				markdowns.GET("/get-markdown/:id", func(c *gin.Context) {
					delivery.GetMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/get-all-markdowns
				markdowns.GET("/get-all-markdowns", func(c *gin.Context) {
					delivery.GetAllMarkdowns(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/delete-markdown
				markdowns.PUT("/delete-markdown/:id", func(c *gin.Context) {
					delivery.DeleteMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/update-markdown
				markdowns.PUT("/update-markdown", func(c *gin.Context) {
					delivery.UpdateMarkdown(a.repository, c)
				})
			}
		}
	}

	// TODO: Add contributor routes

	// contributor := router.Group("/contribute")
	// {

	// }

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
