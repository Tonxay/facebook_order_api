package routers_part

import (
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupCustomersRoutesPart(route fiber.Router) {
	route.Get("/all", service.GetFacebookAllCustomers)
	route.Get("/:facebook_id", service.GetFacebookProfile)

}
