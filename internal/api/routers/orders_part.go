package routers_part

import (
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupOrdersRoutesPart(route fiber.Router) {
	route.Post("/create", service.CreateOrder)

}
