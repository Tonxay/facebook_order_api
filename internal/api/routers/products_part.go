package routers_part

import (
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutesPart(route fiber.Router) {
	route.Post("/create/category", service.CreateCategorie)
	route.Post("/create/product", service.CreateProduct)
	route.Post("/create/product-detail", service.CreateProductDetail)
	route.Post("/create/product-stock", service.CreateStockProductDetail)

}
