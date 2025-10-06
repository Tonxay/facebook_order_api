package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupExpenseRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/", service.GetFullExpenseDetails)
	route.Post("/batch", service.CreateExpensesBatch)
	route.Get("/categorys", service.GetExpensesCategory)
	route.Get("/suppliers", service.GetExpensesSuppliers)
	route.Get("/products", service.GetExpensesProducts)

}
