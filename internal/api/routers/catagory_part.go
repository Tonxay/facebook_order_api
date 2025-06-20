package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupCatagorysRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Post("create/category", service.CreateCategorie)
	route.Get("/all", service.GetCategory)

}
