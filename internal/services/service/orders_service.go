package service

import (
	"fmt"

	cons "go-api/internal/config/constant"
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/services/db_service"

	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetOrder(c *fiber.Ctx) error {
	var err error

	var data []*custommodel.OrderReponse

	var req = request.StatusOrderRequest{}

	// Parse JSON body into struct
	if err := middleware.ParseAndValidateBody(c, &req); err != nil {
		return fiber.NewError(400, err.Error())
	}

	db := gormpkg.GetDB()

	data, err = dbservice.GetOrders(db, req.Statuses, req.IsCancell)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(200).JSON(presenters.ResponseSuccess(data))
}

func CreateOrder(c *fiber.Ctx) error {

	userID, ok := middleware.GetUserID(c)

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
	if req.FacebookID != "" {
		_, err = dbservice.UpdateColumnsCustomer(db, req.FacebookID, int32(req.Gender), req.Tel)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	number := middleware.GenerateOrderNumber()

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	order := models.Order{
		ID:            orderID,
		OrderName:     req.FullName,
		Status:        cons.Ordered,
		CustomerID:    req.FacebookID,
		OrderNo:       number,
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
	orderTimeLine := models.OrderTimeLine{
		OrderStatus: cons.Ordered,
		UserID:      userID,
		OrderID:     orderID,
	}

	err = dbservice.CreateOrder(db, &order)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateOrderTimeLine(db, &orderTimeLine)
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

	respones, err1 := dbservice.GetOrder(db, orderID)

	if err1 != nil {
		return fiber.NewError(http.StatusInternalServerError, err1.Error())
	}

	db.Commit()
	webhookURL := "https://discord.com/api/webhooks/1382990950107971696/gPuNdiZ_0YrxWczQKTtTOxndUkukqyrtlyxg7T63zCcFLgO4JQZzAESuAaKOQEb8QcOy"
	message := fmt.Sprintf(`
ວັນທີ່: %s 
ລູກຄ້າ: %s 
ລະຫັດ: %s
ຈັດສົ່ງໂດຍ: %s
ເບີໂທ: %d
ທີ່ຢູ່: ເເຂວງ %s ເມືອງ %s ສາຂາ %s
ຈາກ: %s
`,
		respones.CreatedAt,
		respones.OrderName,
		respones.OrderNo,
		respones.Shipping.Name,
		order.Tel,
		respones.Province,
		respones.District,
		order.CustomAddress,
		respones.PageName,
	)

	SendDiscordWebhook(webhookURL, message)

	return c.Status(200).JSON(presenters.ResponseSuccess(respones))
}

func UpdateStatusOrder(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	var err error

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}

	status := c.Query("status")
	orderId := c.Query("order_id")
	orderNo := c.Query("order_no")
	oldStatus, ok := cons.OrderStatusTransitions[status]

	if !ok {
		return fiber.NewError(400, "not found status")
	}
	db := gormpkg.GetDB().Begin()

	defer func() {
		db.Rollback()
	}()

	order, err := dbservice.UpdateStatusOrder(db, orderId, status, oldStatus, orderNo)

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	err = dbservice.CreateOrderTimeLine(db, &models.OrderTimeLine{
		UserID:      userID,
		OrderStatus: status,
		OrderID:     order.ID,
	})

	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	db.Commit()

	return c.Status(200).JSON(presenters.ResponseSuccess("update status success"))
}

func CancellOrder(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	var err error

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}

	status := c.Query("status")
	orderId := c.Query("order_id")

	if status != cons.OrderCancelled {
		return fiber.NewError(400, "not found status")
	}
	db := gormpkg.GetDB().Begin()

	defer func() {
		db.Rollback()
	}()

	orderDetail, err1 := dbservice.GetOrderDetails(db, orderId)
	if err1 != nil {
		return fiber.NewError(500, err1.Error())
	}
	err = dbservice.UpdateIsCancelOrder(db, orderId)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if len(orderDetail) <= 0 {
		return fiber.NewError(400, "status is not change")
	}

	var stockProductDetail []*models.StockProductDetail

	for _, item := range orderDetail {
		println(item.ProductDetailID)

		stockProductDetail = append(stockProductDetail, &models.StockProductDetail{
			Quantity:        item.Quantity,
			Remaining:       item.Quantity,
			UserID:          userID,
			Status:          "active",
			ProductDetailID: item.ProductDetailID,
			SizeID:          item.SizeID,
		})

	}

	err2 := dbservice.CreateStockProductDetailForOrder(db, stockProductDetail, c.Context())
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create stock",
		})
	}
	err = dbservice.CreateOrderTimeLine(db, &models.OrderTimeLine{
		UserID:      userID,
		OrderStatus: cons.OrderCancelled,
		OrderID:     orderId,
	})
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	db.Commit()
	return c.Status(200).JSON(presenters.ResponseSuccess("update status success"))
}
