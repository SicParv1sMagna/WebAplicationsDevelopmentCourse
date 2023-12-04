package app

import (
	"log"
	"project/docs"
	"project/internal/http/delivery"
	"project/internal/pkg/roles"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *Application) StartServer() {
	log.Println("Server starting")

	docs.SwaggerInfo.Title = "Notek Rest-API"
	docs.SwaggerInfo.Description = "Notek Server implementation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	user := router.Group("/user")
	{
		user.POST("/register", a.delivery.RegisterUser)
		user.POST("/login", a.delivery.LoginUser)
	}

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/me", a.OnAuthCheck(roles.User, roles.Moderator, roles.Admin), a.delivery.GetMe)
			user.DELETE("/delete", delivery.DeleteUser)
			user.PUT("/edit-info", delivery.UpdateUserInfo)
		}
		notes := api.Group("/notes")
		{
			markdowns := notes.Group("/markdown")
			{
				markdowns.POST("/create", a.OnAuthCheck(roles.User, roles.Moderator, roles.Admin), a.delivery.CreateMarkdown)
				markdowns.GET("/", a.delivery.GetAllMarkdowns)
				markdowns.GET("/:id", a.delivery.GetMarkdown)
				markdowns.DELETE("/:id", a.delivery.DeleteMarkdown)
				markdowns.PUT("/", a.delivery.UpdateMarkdown)
				markdowns.POST("/add-md-to-contributor/:markdown_id/:contributor_id", a.delivery.AddMarkdownToContributor)
				markdowns.POST("/:id/image", a.delivery.AddMarkdownIcon)
			}
		}
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

}
