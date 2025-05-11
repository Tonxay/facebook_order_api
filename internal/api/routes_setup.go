package api

import (
	"github.com/gofiber/fiber/v2"
	routers_part "github.com/yourusername/go-api/internal/api/routers"
)

func SetupRoutes(app *fiber.App) {
	conversation := app.Group("/conversations")
	routers_part.SetupConversationsRoutesPart(conversation)
	webhook := app.Group("/webhook")
	routers_part.SetupWebhookRoutesPart(webhook)

}

func SetupWebsocketRoutes(app *fiber.App) {
	routers_part.SetupWebSocketRoutesPart(app)
}
