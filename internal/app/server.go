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
		//	http://localhost:8080/user/register		ГОТОВО
		user.POST("/register", func(c *gin.Context) {
			delivery.RegisterUser(a.repository, c)
		})

		//	http://localhost:8080/user/login		ГОТОВО
		user.POST("/login", func(c *gin.Context) {
			delivery.AuthUser(a.repository, c)
		})

	}

	//	http://localhost:8080/api
	api := router.Group("/api")
	{
		api.Use(jwttoken.CheckJWTToken())
		//	http://localhost:8080/api/user		ГОТОВО
		user := api.Group("/user")
		{
			//	http://localhost:8080/api/user/delete	ГОТОВО
			user.DELETE("/delete", delivery.DeleteUser)

			//	http://localhost:8080/api/user/edit-info	ГОТОВО
			user.PUT("/edit-info", delivery.UpdateUserInfo)
		}

		//	http://localhost:8080/api/notes
		notes := api.Group("/notes")
		{
			//	http://localhost:8080/api/notes/markdown
			markdowns := notes.Group("/markdown")
			{
				//	http://localhost:8080/api/notes/markdown/create-markdown		ГОТОВО
				markdowns.POST("/create-markdown", func(c *gin.Context) {
					delivery.CreateMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/get-markdown	ГОТОВО
				markdowns.GET("/get-markdown/:id", func(c *gin.Context) {
					delivery.GetMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/get-all-markdowns		ГОТОВО
				markdowns.GET("/get-all-markdowns", func(c *gin.Context) {
					delivery.GetAllMarkdowns(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/delete-markdown		ГОТОВО
				markdowns.DELETE("/delete-markdown/:id", func(c *gin.Context) {
					delivery.DeleteMarkdown(a.repository, c)
				})

				//	http://localhost:8080/api/notes/markdown/update-markdown		ГОТОВО
				markdowns.PUT("/update-markdown", func(c *gin.Context) {
					delivery.UpdateMarkdown(a.repository, c)
				})
			}
		}

		contributor := api.Group("/contributor")
		{
			//	http://localhost:8080/api/contributor/get-contributor/:id	ГОТОВО
			contributor.GET("/get-contributor/:id", func(c *gin.Context) {
				delivery.GetContributor(a.repository, c)
			})

			//	http://localhost:8080/api/contributor/get-all-contributors
			contributor.GET("/get-all-contributors/:id", func(c *gin.Context) {
				delivery.GetAllContributorsFromMarkdown(a.repository, c)
			})

			//	http://localhost:8080/api/contributor/delete-contributor
			contributor.DELETE("/delete-contributor", func(c *gin.Context) {
				delivery.DeleteContributorFromMd(a.repository, c)
			})

			//	http://localhost:8080/api/contibutor/update-contributor
			contributor.PUT("/update-contributor", func(c *gin.Context) {
				delivery.UpdateContributorAccess(a.repository, c)
			})
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
