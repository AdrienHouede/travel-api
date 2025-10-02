package main

import (
	"fmt"
	"log"
	"time"

	"my-gin-project/src/controllers"
	"my-gin-project/src/models"
	"my-gin-project/src/routes"

	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	_ "my-gin-project/src/docs"
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
	// Initialisation de Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://2f1167ff3d20366cfa3695b14e6cb581@o4510114747121664.ingest.de.sentry.io/4510114754592848",
		TracesSampleRate: 1.0, // pour activer le tracing, ajuster selon tes besoins
	})
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	defer sentry.Flush(2 * time.Second) // s'assure que les événements sont envoyés avant la fin du programme

	db, err := models.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Créer le controller avec la DB
	chatController := &controllers.Controller{
		DB: db,
	}

	r := gin.Default()

	// Ajout du middleware Sentry
	r.Use(sentrygin.New(sentrygin.Options{}))

	r.GET("/panic", func(c *gin.Context) {
		defer sentry.Recover() // capture un panic
		panic("Something went wrong!")
	})

	r.GET("/custom-error", func(c *gin.Context) {
		err := fmt.Errorf("une erreur personnalisée")
		sentry.CaptureException(err)
		c.JSON(500, gin.H{"error": err.Error()})
	})

	// Configuration des routes
	routes.SetupRoutes(r, chatController)

	// Lancement du serveur
	r.Run(":8080")
}
