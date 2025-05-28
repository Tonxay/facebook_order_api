package auth

import (
	"go-api/internal/config/middleware"
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthenRoutesPart(route fiber.Router) {
	route.Post("/login", service.Login)
	route.Post("/refresh", service.Refresh)
	route.Post("/register", middleware.JWTProtected, service.Register)
}
