package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetProductforProvinceSevice(c *fiber.Ctx) error {

	productID := c.Query("product_id")

	err := middleware.ValidateUuId(productID)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	data, err := dbservice.GetProductforProvince(gormpkg.GetDB(), productID)

	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.JSON(presenters.ResponseSuccess(data))

}
