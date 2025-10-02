package routes

import (
	"my-gin-project/src/controllers"

	"github.com/gin-gonic/gin"

	_ "my-gin-project/src/docs" // import Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, ctrl *controllers.Controller) {

	// Routes publiques
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)

	// Ajout des chatbot
	router.POST("/chat", ctrl.Chat)
	router.POST("/chat-ai", ctrl.ChatAI)

	// Routes protégées
	authorized := router.Group("/")
	authorized.Use(controllers.AuthMiddleware())
	{
		authorized.GET("/items", ctrl.GetItems)
		authorized.GET("/items/:id", ctrl.GetItemByID)
		authorized.POST("/items", ctrl.CreateItem)
		authorized.PUT("/items/:id", ctrl.UpdateItem)
		authorized.DELETE("/items/:id", ctrl.DeleteItem)
	}

	// Route Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}
