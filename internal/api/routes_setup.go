package api

import (
	"github.com/gofiber/fiber/v2"
	routers_part "github.com/yourusername/go-api/internal/api/routers"
)

func SetupRoutes(app *fiber.App) {
	webhook := app.Group("/webhook")
	routers_part.SetupRoutesPart(webhook)
}
