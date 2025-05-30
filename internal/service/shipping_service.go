package service

import (
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	dbservice "go-api/internal/service/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetShipping(c *fiber.Ctx) error {
	data, err := dbservice.GetShipping(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get shipping",
		})
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}
