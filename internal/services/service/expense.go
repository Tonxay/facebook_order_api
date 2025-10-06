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

func GetFullExpenseDetails(c *fiber.Ctx) error {

	var err error
	filter := custommodel.ExpenseFilter{}
	err = c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	filter.StartDate, filter.EndDate = middleware.SetDefaultDateRangeMonthIfEmpty(filter.StartDate, filter.EndDate)

	data, err := dbservice.GetFullExpenseDetails(gormpkg.GetDB(), filter)
	if err != nil {
		return fiber.NewError(500, "server error")
	}
	return c.JSON(presenters.ResponseSuccess(data))

}

// CreateExpenseWithProduct creates a new expense, including a potential product link.
func CreateExpenseWithProduct(c *fiber.Ctx) error {
	var input []*custommodel.ExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := dbservice.CreateExpensesInBatch(gormpkg.GetDB(), input)
	if err != nil {
		return fiber.NewError(500, "server error")
	}

	return c.Status(fiber.StatusCreated).JSON(data)
}

// CreateExpensesBatch handles creating multiple expenses.
func CreateExpensesBatch(c *fiber.Ctx) error {
	var inputs []*custommodel.ExpenseInput // Use a slice of the unified input struct
	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if len(inputs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body must be a non-empty array"})
	}

	data, err := dbservice.CreateExpensesInBatch(gormpkg.GetDB(), inputs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully created %d expenses", len(data)),
		"data":    data,
	})
}

func GetExpensesCategory(c *fiber.Ctx) error {

	data, err := dbservice.GetFullExpenseCategory(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(presenters.ResponseSuccess(data))
}

func GetExpensesSuppliers(c *fiber.Ctx) error {

	data, err := dbservice.GetFullExpenseSuppliers(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(presenters.ResponseSuccess(data))

}

func GetExpensesProducts(c *fiber.Ctx) error {

	data, err := dbservice.GetFullExpenseProducts(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(presenters.ResponseSuccess(data))
}
