package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/services/db_service"

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

	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(requestData))
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

	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(requestData))
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

	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(requestData))
}
func CreateStockProductDetail(c *fiber.Ctx) error {
	id, _ := middleware.GetUserID(c)
	var requestData request.StockProductDetail

	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	data := models.StockProductDetail{
		UserID:          id,
		ProductDetailID: requestData.ProductDetailID,
		Quantity:        requestData.Quantity,
		SizeID:          requestData.SizeID,
		Remaining:       requestData.Quantity,
	}
	db := gormpkg.GetDB()

	err := dbservice.CreateStockProductDetail(db, &data, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create stock",
		})
	}
	// Return the created category

	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
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

func GetProductInStcok(c *fiber.Ctx) error {
	data, err := dbservice.GetProducts(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get products",
		})
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}
