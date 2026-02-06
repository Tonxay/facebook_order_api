package routers_part

import (
	"go-api/internal/services/service"

	"github.com/gofiber/fiber/v2"
)

func SetupTictokRoutesPart(route fiber.Router) {
	// route.Use(middleware.JWTProtected)
	route.Get("/download", service.TictokVideoService)

	// route.Get("/pinduoduo/download", service.PinduoduoService)
}
