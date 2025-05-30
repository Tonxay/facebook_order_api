package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupShippingRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/", service.GetShipping)

}
