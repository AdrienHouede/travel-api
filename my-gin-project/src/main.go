package main

import (
	"log"
	"my-gin-project/src/models"
	"my-gin-project/src/routes"

	"github.com/gin-gonic/gin"

	_ "my-gin-project/src/docs" // Swagger docs
)

// @title My Gin API
// @version 1.0
// @description API pour gérer des items et l'authentification JWT.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialisation DB
	if err := models.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// Configuration des routes
	routes.SetupRoutes(r)

	// Lancement du serveur
	r.Run(":8080")
}
