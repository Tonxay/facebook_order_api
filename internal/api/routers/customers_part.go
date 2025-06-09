package routers_part

import (
	"go-api/internal/config/middleware"
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupCustomersRoutesPart(route fiber.Router) {
	route.Use(middleware.JWTProtected)
	route.Get("/all", service.GetFacebookAllCustomers)
	route.Get("/:facebook_id", service.GetFacebookProfile)

}
