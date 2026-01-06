package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupBillRoutesPart(route fiber.Router) {
	// route.Get("/send", service.Uploadimage)

	route.Get("/anousith", middleware.JWTProtected, service.GetOrderbillInAnousith)
	route.Post("/screenshot", service.ScapingImage)
	route.Get("proxy-anousith", service.ProxyAnousith)
	route.Get("/:tracking_number", service.Getbill)

}
