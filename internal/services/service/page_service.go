package service

import (
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetPages(c *fiber.Ctx) error {
	pagetype := c.Query("page_type")
	data, err := dbservice.GetPages(gormpkg.GetDB(), pagetype)
	if err != nil {
		return fiber.NewError(500, "server error")
	}
	return c.JSON(presenters.ResponseSuccess(data))
}
