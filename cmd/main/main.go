package main

import (
	"errors"
	"go-api/internal/api"
	"go-api/internal/config/middleware"
	gormpkg "go-api/internal/pkg"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// API Routes
	// api := app.Group(os.Getenv("API_PREFIX"))

	myConfig := fiber.Config{
		// DisableStartupMessage: true,
		// AppName: apiName,
		// Override default body limit to 50MB
		BodyLimit: 50 * 1024 * 1024,
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error response
			err = ctx.Status(code).JSON(fiber.Map{
				"timestamp": time.Now().Format("2006-01-02-15-04-05"),
				"status":    0,
				"items":     nil,
				"error":     err.Error(),
			})

			// Return from handler
			return err
		},
	}

	app := fiber.New(myConfig)
	middleware.Init()
	api.SetupRoutes(app)
	if err := gormpkg.Init("webhook"); err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	log.Println("run app ...")
	api.SetupWebsocketRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
