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
		return fiber.NewError(500, "failed to get shipping")
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "failed to get shipping",
		// })
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}

func GetProvices(c *fiber.Ctx) error {
	data, err := dbservice.GetProvince(gormpkg.GetDB())
	if err != nil {
		return fiber.NewError(500, "failed to get provices")
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}
func GetDistricts(c *fiber.Ctx) error {
	provice_id := c.Params("provice_id")
	data, err := dbservice.GetDistrict(gormpkg.GetDB(), provice_id)
	if err != nil {
		return fiber.NewError(500, "failed to get districts")
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}
