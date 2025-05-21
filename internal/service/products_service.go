package service

import (
	"go-api/internal/pkg/models"
	dbservice "go-api/internal/service/db_service"

	"github.com/gofiber/fiber/v2"
)

func CreateCategorie(c *fiber.Ctx) error {
	var requestData models.Category

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	err := dbservice.CreateCategory(&requestData, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create category",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": requestData,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	var requestData models.Product

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	err := dbservice.CreateProduct(&requestData, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create product",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": requestData,
	})
}

func CreateProductDetail(c *fiber.Ctx) error {
	var requestData models.ProductDetail

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	err := dbservice.CreateProductDetail(&requestData, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create product",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": requestData,
	})
}
func CreateStockProductDetail(c *fiber.Ctx) error {
	var requestData models.StockProductDetail

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	err := dbservice.CreateStockProductDetail(&requestData, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create sock",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": requestData,
	})
}
