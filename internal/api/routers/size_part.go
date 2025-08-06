package routers_part

import (
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupSizeRoutesPart(route fiber.Router) {
	route.Post("/create", service.CreateSize)
}
