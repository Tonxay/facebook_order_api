package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/products-province", service.GetProductforProvinceSevice)
	route.Get("/order-summary", service.GetOrderSummary)
	route.Get("/product-sales-day", service.GetProductSalesByDay)
	route.Get("/product-order-count", service.GetProductOrderCount)
	route.Get("/product-sales", service.GetProductSales)

}
