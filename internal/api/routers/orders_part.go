package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupOrdersRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Post("/create", service.CreateOrder)
	route.Post("/all", service.GetOrder)

	route.Put("/update-status", service.UpdateStatusOrder)
	route.Put("/cancell", service.CancellOrder)

}
