package service

import (
	"fmt"
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetProductforProvinceSevice(c *fiber.Ctx) error {
	var err error

	filter := custommodel.FilterDasboard{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	err = middleware.ValidateUuId(filter.ProductID)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeMonthIfEmpty(filter.StartDate, filter.EndDate)

	data, err := dbservice.GetProductforProvince(gormpkg.GetDB(), filter)

	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.JSON(presenters.ResponseSuccess(data))

}

func GetOrderSummary(c *fiber.Ctx) error {
	var err error
	filter := custommodel.FilterDasboard{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	// filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeIfEmpty(filter.StartDate, filter.EndDate)

	data, err := dbservice.GetOrderSummary(gormpkg.GetDB(), filter)

	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.JSON(presenters.ResponseSuccess(data))

}

func GetProductSalesByDay(c *fiber.Ctx) error {

	var err error
	filter := custommodel.FilterDasboard{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeIfEmpty(filter.StartDate, filter.EndDate)

	data, err := dbservice.GetProductSalesByDay(gormpkg.GetDB(), filter)

	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.JSON(presenters.ResponseSuccess(data))

}
func GetProductOrderCount(c *fiber.Ctx) error {
	var err error
	filter := custommodel.FilterDasboard{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeIfEmpty(filter.StartDate, filter.EndDate)

	data, err := dbservice.GetProductOrderCount(gormpkg.GetDB(), filter)

	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.JSON(presenters.ResponseSuccess(data))

}

func GetProductSales(c *fiber.Ctx) error {

	var err error
	filter := custommodel.FilterDasboard{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}
	fmt.Println(filter.StartDate)
	filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeIfEmpty(filter.StartDate, filter.EndDate)
	fmt.Println(filter.StartDate)
	data, err := dbservice.GetProductSales(gormpkg.GetDB(), filter)
	if err != nil {
		return fiber.NewError(500, "server error")
	}
	return c.JSON(presenters.ResponseSuccess(data))

}
