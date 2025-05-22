package service

import (
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/custom_model/request"
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
	var requestData request.StockProductDetail

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	data := models.StockProductDetail{
		UserID:          "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77",
		ProductDetailID: requestData.ProductDetailID,
		Quantity:        requestData.Quantity,
		SizeID:          requestData.SizeID,
		Remaining:       requestData.Quantity,
	}

	err := dbservice.CreateStockProductDetail(&data, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create stock",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": data,
	})
}

func GetStockProductDetailForID(c *fiber.Ctx) error {
	id := c.Query("product_item_detail_id")
	size_id := c.Query("size_id")
	data, err := dbservice.GetProductDetailsForID(gormpkg.GetDB(), id, size_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get stock",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(data)
}

func CreateProductSize(c *fiber.Ctx) error {

	var requestData models.Size

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}
	err := dbservice.CreateProductSize(&requestData, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create size",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(requestData)
}
