package service

import (
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetCategory(c *fiber.Ctx) error {
	data, err := dbservice.GetCategory(gormpkg.GetDB())
	if err != nil {
		return fiber.NewError(500, "server error")
	}
	return c.JSON(presenters.ResponseSuccess(data))
}
