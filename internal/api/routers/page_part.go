package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupPagesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/all", service.GetPages)
}
