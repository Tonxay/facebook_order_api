package routers_part

import (
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupWebhookRoutesPart(route fiber.Router) {
	route.Get("/", service.VerifyWebhook)
	route.Post("/", service.HandleWebhook)

}
