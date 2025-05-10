package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go-api/internal/api"
	gormpkg "github.com/yourusername/go-api/internal/pkg"
)

func main() {
	if err := gormpkg.Init("api"); err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	app := fiber.New()
	// API Routes
	// api := app.Group(os.Getenv("API_PREFIX"))
	api.SetupRoutes(app)
	api.SetupWebsocketRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
