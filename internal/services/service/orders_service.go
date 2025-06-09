package service

import (
	"fmt"
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/services/db_service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateOrder(c *fiber.Ctx) error {
	user := c.Locals("user_id") // returns interface{}
	userID, ok := user.(string) // type assert to string (or int, etc.)

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}
	var req custommodel.OrderRequest

	var err error

	// Parse JSON body into struct
	if err := middleware.ParseAndValidateBody(c, &req); err != nil {
		return fiber.NewError(400, err.Error())
	}

	db := gormpkg.GetDB().Begin()
	defer func() {
		db.Rollback()
	}()
	// Step 1: Group items by product_id
	temp := make(map[string][]custommodel.GroupedItem)
	total_quantity := make(map[string]int32)
	TotalPrice := make(map[string]float64)
	price := make(map[string]int)
	var data custommodel.ProductOrderDetails

	for _, item := range req.Items {

		data, err = dbservice.GetProductDetailsByIDSizdID(db, item.ProductDetailID, item.SizeID, item.ProductID)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		if data.Remaining < int32(item.Quantity) {
			data, err = dbservice.GetProductByIDSizdID(db, item.ProductDetailID, item.SizeID, item.ProductID)
			if err != nil {
				return fiber.NewError(400, err.Error())
			}

			return fiber.NewError(
				550,
				fmt.Sprintf("%s %s %s ຍັງເຫລືອເເຕ່ %s ໂຕ",
					data.Name,
					data.ColorName,
					data.Size,
					strconv.Itoa(int(data.Remaining)),
				),
			)
		}
		grouped := custommodel.GroupedItem{
			ProductDetailID: item.ProductDetailID,
			SizeID:          item.SizeID,
			ProductDetails:  data,
			Quantity:        item.Quantity,
		}
		temp[item.ProductID] = append(temp[item.ProductID], grouped)
		total_quantity[item.ProductID] += item.Quantity
		price[item.ProductID] = int(data.Price)
		TotalPrice[item.ProductID] += float64((data.Price * float64(item.Quantity)))
	}

	var discount_only_product float32
	// Step 2: Build final output
	var groupedResult []custommodel.GroupedByProduct
	for productID, groupedItems := range temp {
		promotions, _ := dbservice.GetPromotion(db, productID)
		total_quantity := total_quantity[productID]
		var pot models.Promotion
		for _, promtion := range promotions {

			if total_quantity == promtion.Quentity {
				pot = promtion
				discount_only_product = float32(total_quantity/promtion.Quentity) * promtion.Discount
			}
			if total_quantity > promtion.Quentity {
				pot = promotions[len(promotions)-1]
				discount_only_product = float32(total_quantity/pot.Quentity) * pot.Discount
			}

		}

		groupedResult = append(groupedResult, custommodel.GroupedByProduct{
			ProductID: productID,

			PromotionID:         pot.ID,
			DiscountOnlyProduct: discount_only_product,
			Promotion:           promotions,
			TotalPrice:          TotalPrice[productID],
			TotalQuantities:     total_quantity,
			Items:               groupedItems,
		})
	}
	var orderDetail []*models.OrderDetail
	var orderDiscount []*models.OrderDiscount
	var orderStockDetails []*models.OrderStockDetail

	orderID := uuid.New().String()
	var totalPrice float64

	for _, data := range groupedResult {

		for _, item := range data.Items {
			product, err := dbservice.GetProductDetailsForID(db, item.ProductDetailID, item.SizeID)
			if err != nil {
				return fiber.NewError(http.StatusBadRequest, "not found product")
			}
			productQuantity := int32(item.Quantity)
			for _, stock := range product.StockProducts {

				if productQuantity >= 1 {

					if stock.Remaining >= productQuantity {
						stock.Remaining -= productQuantity

						orderStockDetails = append(orderStockDetails, &models.OrderStockDetail{
							OrderID:              orderID,
							StockProductDetailID: stock.ID,
							Quatity:              productQuantity,
						})

						if stock.Remaining >= 1 {
							productQuantity = 0
							err = dbservice.UpdateStockProductDetail(db, stock.ID, stock.Remaining, "active")
						} else {
							err = dbservice.UpdateStockProductDetail(db, stock.ID, stock.Remaining, "out_stock")
						}

						if err != nil {
							return fiber.NewError(http.StatusInternalServerError, "server create order details error")
						}

					}

				}

			}
			orderDetail = append(orderDetail, &models.OrderDetail{
				OrderID:         orderID,
				ProductDetailID: item.ProductDetailID,
				Quantity:        int32(item.Quantity),
				SizeID:          item.SizeID,
				UnitPrice:       float64(product.Price),
				TotalPrice:      float64((product.Price * item.Quantity)),
			})

		}

		totalPrice += data.TotalPrice

		if data.DiscountOnlyProduct > 0 {
			orderDiscount = append(orderDiscount, &models.OrderDiscount{
				OrderID:       orderID,
				ProductID:     data.ProductID,
				TotalDiscount: float64(data.DiscountOnlyProduct),
				DiscountID:    data.PromotionID,
			})
		}
	}

	order := models.Order{
		ID:            orderID,
		OrderName:     req.FullName,
		Status:        "pending",
		CustomerID:    req.FacebookID,
		OrderNo:       middleware.GenerateOrderNumber(),
		Tel:           req.Tel,
		ShippingID:    req.ShippingID,
		CustomAddress: req.CustomAddress,
		UserID:        userID,
		DistrictID:    req.DistrictID,
		TotalPrice:    totalPrice,
		FreeShipping:  req.FreeShipping,
		Cod:           req.Cod,
		Discount:      req.Discount,
		Platform:      req.PlatForm,
	}

	err = dbservice.CreateOrder(db, &order)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateOrderDiscounts(db, orderDiscount, c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateRoderStockDetails(db, orderStockDetails, c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	db.Commit()

	return c.Status(200).JSON(presenters.ResponseSuccess(fiber.Map{
		"order": groupedResult,
		// "items": groupedResult,
	}))
}
