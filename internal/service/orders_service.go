package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/service/db_service"
	"net/http"

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
	var totalPrice int32

	for _, item := range req.OrderDetails {
		product, err := dbservice.GetProductDetailsForID(db, item.ProductDetailID, item.SizeID)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "not found product")

		}

		if item.Quantity > product.Quantity {
			return fiber.NewError(http.StatusBadRequest, "quantity than product")

		}

		productQuantity := item.Quantity
		for _, stock := range product.StockProducts {

			if stock.Remaining <= 0 {
				continue
			}
			if productQuantity >= 1 {

				if stock.Remaining >= productQuantity {
					stock.Remaining -= productQuantity
					productQuantity = 0
					err = dbservice.UpdateStockProductDetail(stock.ID, stock.Remaining, "active", c.Context())
					if err != nil {
						return fiber.NewError(http.StatusInternalServerError, "server create order details error")
					}
				} else {
					productQuantity -= stock.Remaining
					stock.Remaining = 0
					err = dbservice.UpdateStockProductDetail(stock.ID, stock.Remaining, "out_stock", c.Context())
					if err != nil {
						return fiber.NewError(http.StatusInternalServerError, "server create order details error")

					}
				}

			}

		}

		orderDetail = append(orderDetail, &models.OrderDetail{
			OrderID:         orderID,
			ProductDetailID: item.ProductDetailID,
			Quantity:        item.Quantity,
			SizeID:          item.SizeID,
			UnitPrice:       float64(product.Price),
			TotalPrice:      float64(item.Quantity * product.Price),
		})

		totalPrice = (item.Quantity * product.Price) + totalPrice

	}

	order := models.Order{
		ID:            orderID,
		Status:        "pending",
		CustomerID:    req.CustomerID,
		PackagePrice:  0,
		OrderNo:       middleware.GenerateOrderNumber(),
		Tel:           req.Tel,
		CustomAddress: req.CustomAddress,
		UserID:        "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77",
		DistrictID:    req.DistrictID,
		TotalPrice:    totalPrice,
	}

	err = dbservice.CreateOrder(db, &order, c.Context())
	if err != nil {

		return fiber.NewError(http.StatusInternalServerError, "server create order details error")

	}

	err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
	if err != nil {

		return fiber.NewError(http.StatusInternalServerError, "server create order details error")
	}
	db.Commit()

	return c.Status(200).JSON(presenters.ResponseSuccess(fiber.Map{
		"order": order,
		"items": orderDetail,
	}))

}
