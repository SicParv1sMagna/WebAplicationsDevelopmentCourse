package app

import (
	"log"
	"project/internal/delivery"
	jwttoken "project/internal/middleware/jwt"

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
			// УСЛУГИ
			//	http://localhost:8080/api/notes/markdown
			markdowns := notes.Group("/markdown")
			{
				// ДОБАВЛЕНИЕ УСЛУГИ
				//	http://localhost:8080/api/notes/markdown/create-markdown		ГОТОВО
				markdowns.POST("/create", func(c *gin.Context) {
					delivery.CreateMarkdown(a.repository, c)
				})

				// ПОЛУЧЕНИЕ УСЛУГ
				//	http://localhost:8080/api/notes/markdown/:id	ГОТОВО
				markdowns.GET("/:id", func(c *gin.Context) {
					delivery.GetMarkdown(a.repository, c)
				})

				// СПИСОК УСЛУГ
				//	http://localhost:8080/api/notes/markdown/		ГОТОВО
				markdowns.GET("/", func(c *gin.Context) {
					delivery.GetAllMarkdowns(a.repository, c)
				})

				// УДАЛЕНИЕ УСЛУГИ
				//	http://localhost:8080/api/notes/markdown/delete-markdown		ГОТОВО
				markdowns.DELETE("/:id", func(c *gin.Context) {
					delivery.DeleteMarkdown(a.repository, c)
				})

				// РЕДАКТИРОВАНИЕ УСЛУГИ
				//	http://localhost:8080/api/notes/markdown/update-markdown		ГОТОВО
				markdowns.PUT("/", func(c *gin.Context) {
					delivery.UpdateMarkdown(a.repository, c)
				})

				// ДОБАВЛЕНИЕ УСЛУГИ В ПОСЛЕДНЮЮ ЗАЯВКУ
				markdowns.POST("/:id/contributor", func(c *gin.Context) {
					delivery.AddMarkdownToContributor(a.repository, c)
				})

				markdowns.DELETE("/:id/contributor/delete", func(c *gin.Context) {
					delivery.DeleteContributorFromMd(a.repository, c)
				})

				markdowns.POST("/:id/image", func(c *gin.Context) {
					delivery.AddMarkdownIcon(a.repository, c)
				})
			}
		}

		// ЗАЯВКИ
		contributor := api.Group("/contributor")
		{
			// ПОЛУЧЕНИЕ ЗАЯВКИ
			//	http://localhost:8080/api/contributor/get-contributor/:id	ГОТОВО
			contributor.GET("/:id", func(c *gin.Context) {
				delivery.GetContributor(a.repository, c)
			})

			// СПИСОК ЗАЯВОК
			//	http://localhost:8080/api/contributor/get-all-contributors
			contributor.GET("/", func(c *gin.Context) {
				delivery.GetAllContributorsFromMarkdown(a.repository, c)
			})

			contributor.DELETE("/:id/delete")
			// РЕДАКТИРОВАНИЕ ЗАЯВКИ
			//	http://localhost:8080/api/contibutor/update-contributor
			contributor.PUT("/", func(c *gin.Context) {
				delivery.UpdateContributorAccess(a.repository, c)
			})

			contributor.PUT("/:id/status/user")
			contributor.PUT("/:id/status/moderator")
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
