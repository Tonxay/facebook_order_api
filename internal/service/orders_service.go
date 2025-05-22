package service

import (
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/service/db_service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateOrder(c *fiber.Ctx) error {

	var req custommodel.OrderRequest

	var err error

	// Parse JSON body into struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	db := gormpkg.GetDB().Begin()
	defer func() {
		db.Rollback()
	}()

	var orderDetail []*models.OrderDetail

	orderID := uuid.New().String()
	for _, item := range req.OrderDetails {

		product, err := dbservice.GetProductDetailsForID(db, item.ProductDetailID)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found product",
			})

		}

		if item.Quantity > product.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "quantity than product",
			})
		}
		productQuantity := item.Quantity

		for _, stock := range product.StockProducts {
			if productQuantity <= 0 {
				break
			}
			if stock.Remaining <= 0 {
				continue
			}

			if stock.Remaining >= productQuantity {
				stock.Remaining -= productQuantity
				productQuantity = 0
			} else {
				productQuantity -= stock.Remaining
				stock.Remaining = 0
			}
		}

		orderDetail = append(orderDetail, &models.OrderDetail{
			OrderID:         orderID,
			ProductDetailID: item.ProductDetailID,
			Quantity:        item.Quantity,
			UnitPrice:       float64(product.Price),
			TotalPrice:      float64(item.Quantity * product.Price),
		})

	}

	order := models.Order{
		ID:            orderID,
		Status:        "pending",
		CustomerID:    req.CustomerID,
		Tel:           req.Tel,
		CustomAddress: req.CustomAddress,
		UserID:        "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77",
		DistrictID:    req.DistrictID,
		TotalPrice:    200,
	}

	err = dbservice.CreateOrder(db, &order, c.Context())
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server create order error ",
		})

	}

	err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server create order details error",
		})

	}
	db.Commit()

	return c.Status(200).JSON("successfully")

}
