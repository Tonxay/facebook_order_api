package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/services/db_service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateProductManyItem(c *fiber.Ctx) error {
	var requestData request.ProductManyRequest
	var err error
	// Parse JSON request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}
	db := gormpkg.GetDB().Begin()
	defer func() {
		db.Rollback()
	}()

	productID := uuid.New().String()

	product := models.Product{
		ID:         productID,
		Brand:      requestData.Brand,
		Name:       requestData.ProductName,
		CategoryID: requestData.CategoryID,
	}

	productDetails := []*models.ProductDetail{}
	sizes := []*models.Size{}

	for _, productItem := range requestData.ProductDetails {

		productDetailID := uuid.New().String()

		newProductDetail := models.ProductDetail{
			ID:        productDetailID,
			ProductID: productID,
			Color:     productItem.Color,
			ColorName: productItem.ColorName,
			FitType:   productItem.FitType,
			Material:  productItem.Material,
			ImageURL:  productItem.ImageURL,
		}

		productDetails = append(productDetails, &newProductDetail)

		for _, size := range productItem.Sizes {

			newSize := models.Size{
				ProductDetailID: productDetailID,
				Size:            size.Size,
				Price:           size.Price,
			}

			sizes = append(sizes, &newSize)
		}

	}
	err = dbservice.CreateProduct(db, &product, c.Context())
	if err != nil {
		return fiber.NewError(500, "create product error")
	}
	err = dbservice.CreateProductDetailList(db, productDetails, c.Context())
	if err != nil {
		return fiber.NewError(500, "create product details error")
	}
	err = dbservice.CreateProductSizeList(db, sizes, c.Context())
	if err != nil {
		return fiber.NewError(500, "create product details error")
	}

	db.Commit()
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(product))
}

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

	err := dbservice.CreateProduct(gormpkg.GetDB(), &requestData, c.Context())
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
func GetProductsForStock(c *fiber.Ctx) error {
	data, err := dbservice.GetProductsForStock(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get products",
		})
	}
	// Return the created category
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(data))
}
