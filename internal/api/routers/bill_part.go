package routers_part

import (
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupBillRoutesPart(route fiber.Router) {
	route.Get("/send", service.Uploadimage)
	route.Get("/:tracking_number", service.Getbill)

}
