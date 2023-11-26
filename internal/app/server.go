package app

import (
	"log"
	"project/internal/http/delivery"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server starting")

	// Создаем роутинг
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // List of allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Enable credentials (e.g., cookies)
	}))

	//	http://localhost:8080/user
	user := router.Group("/user")
	{
		user.POST("/register", a.delivery.RegisterUser)
		user.POST("/login", a.delivery.LoginUser)
	}

	//	http://localhost:8080/api
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.DELETE("/delete", delivery.DeleteUser)
			user.PUT("/edit-info", delivery.UpdateUserInfo)
		}

		//	http://localhost:8080/api/notes
		notes := api.Group("/notes")
		{
			// УСЛУГИ
			//	http://localhost:8080/api/notes/markdown
			markdowns := notes.Group("/markdown")
			{
				markdowns.POST("/create", a.delivery.CreateMarkdown)
				markdowns.GET("/", a.delivery.GetAllMarkdowns)
				markdowns.GET("/:id", a.delivery.GetMarkdown)
				markdowns.DELETE("/:id", a.delivery.DeleteMarkdown)
				markdowns.PUT("/", a.delivery.UpdateMarkdown)
				markdowns.POST("/add-md-to-contributor/:markdown_id/:contributor_id", a.delivery.AddMarkdownToContributor)
				markdowns.POST("/:id/image", a.delivery.AddMarkdownIcon)
			}
		}

		// ЗАЯВКИ
		contributor := api.Group("/contributor")
		{
			contributor.GET("/:id", a.delivery.GetContributor)
			contributor.GET("/:id/markdown", a.delivery.GetAllContributorsFromMarkdown)
			contributor.GET("/", a.delivery.GetAllContirbutors)
			contributor.DELETE("/delete", a.delivery.DeleteContributorFromMd)
			contributor.PUT("/moderator", a.delivery.UpdateContributorAccessByModerator)
			contributor.PUT("/admin", a.delivery.UpdateContributroAccessByAdmin)
			contributor.PUT("/:id/user", a.delivery.RequestContribution)
			contributor.PUT("/:id", a.delivery.UpdateContributorData)
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
