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

	route.Put("/update-status", service.UpdateOrder)
	// route.Put("/packed", service.GetOrder)
	// route.Put("/shipped", service.GetOrder)
	// route.Put("/customer-bill-notified", service.GetOrder)
	// route.Put("/delivery-complete", service.GetOrder)
	// route.Put("/payment-completed", service.GetOrder)
	// route.Put("/order-cancelled", service.GetOrder)
	// route.Put("/return-to-sender", service.GetOrder)
	// route.Put("/customer-notified", service.GetOrder)

}
