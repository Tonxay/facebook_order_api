package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Post("/create/category", service.CreateCategorie)
	// route.Post("/create/product", service.CreateProduct)
	route.Post("/create/product-detail", service.CreateProductDetail)
	route.Post("/create/product-stock", service.CreateStockProductDetail)
	route.Post("/create/size", service.CreateProductSize)
	route.Post("/create/products", service.CreateProductManyItem)

	route.Get("/product-detail-size", service.GetStockProductDetailForID)
	route.Get("/product-in-stock", service.GetProductInStcok)
	route.Get("/products-manage", service.GetProductsForStock)

}
