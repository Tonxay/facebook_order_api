package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupShippingRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/provices", service.GetProvices)
	route.Get("/districts/:provice_id", service.GetDistricts)
	route.Get("/", service.GetShipping)

}
