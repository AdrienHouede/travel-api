package routes

import (
	"my-gin-project/src/controllers"

	"github.com/gin-gonic/gin"

	_ "my-gin-project/src/docs" // import Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	controller := controllers.Controller{}

	// Routes publiques
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Ajout des chatbot
	router.POST("/chat", controller.Chat)
	router.POST("/chat-ai", controller.ChatAI)

	// Routes protégées
	authorized := router.Group("/")
	authorized.Use(controllers.AuthMiddleware())
	{
		authorized.GET("/items", controller.GetItems)
		authorized.GET("/items/:id", controller.GetItemByID)
		authorized.POST("/items", controller.CreateItem)
		authorized.PUT("/items/:id", controller.UpdateItem)
		authorized.DELETE("/items/:id", controller.DeleteItem)
	}

	// Route Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}
