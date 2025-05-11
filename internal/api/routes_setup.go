package api

import (
	routers_part "go-api/internal/api/routers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	webhook := app.Group("/webhook")
	routers_part.SetupWebhookRoutesPart(webhook)

	conversation := app.Group("/conversations")
	routers_part.SetupConversationsRoutesPart(conversation)
}

func SetupWebsocketRoutes(app *fiber.App) {
	routers_part.SetupWebSocketRoutesPart(app)
}
