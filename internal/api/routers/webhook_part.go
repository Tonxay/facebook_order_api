package routers_part

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go-api/internal/service"
)

func SetupWebhookRoutesPart(route fiber.Router) {
	route.Get("/", service.VerifyWebhook)
	route.Post("/", service.HandleWebhook)

}
