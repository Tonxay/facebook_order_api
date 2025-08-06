package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"

	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
)

func CreateSize(c *fiber.Ctx) error {
	var err error
	var size models.Size

	err = middleware.ParseAndValidateBody(c, &size)

	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	data, err := dbservice.CreateSize(gormpkg.GetDB(), size)

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(presenters.ResponseSuccess(data))

}
