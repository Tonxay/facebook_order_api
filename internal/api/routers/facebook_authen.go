package routers_part

import (
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func FacebookRoutesPart(route fiber.Router) {
	route.Post("/oauth", service.FacebookCallbackHandler)
}
