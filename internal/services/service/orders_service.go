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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetOrder(c *fiber.Ctx) error {
	var err error

	var data []*custommodel.OrderReponseNew

	var req = request.StatusOrderRequest{}

	// Parse JSON body into struct
	if err := middleware.ParseAndValidateBody(c, &req); err != nil {
		return fiber.NewError(400, err.Error())
	}

	db := gormpkg.GetDB()

	data, err = dbservice.GetOrders(db, req)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(200).JSON(presenters.ResponseSuccess(data))
}

// func CreateOrder(c *fiber.Ctx) error {

// 	userID, ok := middleware.GetUserID(c)

// 	if !ok {
// 		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
// 	}
// 	var req custommodel.OrderRequest

// 	var err error

// 	// Parse JSON body into struct
// 	if err := middleware.ParseAndValidateBody(c, &req); err != nil {
// 		return fiber.NewError(400, err.Error())
// 	}

// 	if !middleware.ValidatePhone(req.Tel) {
// 		return fiber.NewError(400, "ເບີໂທລະສັບບໍ່ຖືກຕ້ອງ")
// 	}
// 	print()

// 	db := gormpkg.GetDB().Begin()
// 	defer func() {
// 		db.Rollback()
// 	}()
// 	// Step 1: Group items by product_id
// 	temp := make(map[string][]custommodel.GroupedItem)
// 	total_quantity := make(map[string]int32)
// 	TotalPrice := make(map[string]float64)
// 	price := make(map[string]int)
// 	var data custommodel.ProductOrderDetails

// 	for _, item := range req.Items {

// 		data, err = dbservice.GetProductDetailsByIDSizdID(db, item.ProductDetailID, item.SizeID, item.ProductID)
// 		if err != nil {
// 			return fiber.NewError(400, err.Error())
// 		}

// 		if data.Remaining < int32(item.Quantity) {
// 			data, err = dbservice.GetProductByIDSizdID(db, item.ProductDetailID, item.SizeID, item.ProductID)
// 			if err != nil {
// 				return fiber.NewError(400, err.Error())
// 			}

// 			return fiber.NewError(
// 				550,
// 				fmt.Sprintf("%s %s %s ຍັງເຫລືອເເຕ່ %s ໂຕ",
// 					data.Name,
// 					data.ColorName,
// 					data.Size,
// 					strconv.Itoa(int(data.Remaining)),
// 				),
// 			)
// 		}
// 		grouped := custommodel.GroupedItem{
// 			ProductDetailID: item.ProductDetailID,
// 			SizeID:          item.SizeID,
// 			ProductDetails:  data,
// 			Quantity:        item.Quantity,
// 		}
// 		temp[item.ProductID] = append(temp[item.ProductID], grouped)
// 		total_quantity[item.ProductID] += item.Quantity
// 		price[item.ProductID] = int(data.Price)
// 		TotalPrice[item.ProductID] += float64((data.Price * float64(item.Quantity)))
// 	}

// 	var discount_only_product float32
// 	// Step 2: Build final output
// 	var groupedResult []custommodel.GroupedByProduct
// 	for productID, groupedItems := range temp {
// 		promotions, _ := dbservice.GetPromotion(db, productID)
// 		total_quantity := total_quantity[productID]
// 		var pot models.Promotion
// 		for _, promtion := range promotions {

// 			if total_quantity == promtion.Quentity {
// 				pot = promtion
// 				discount_only_product = float32(total_quantity/promtion.Quentity) * promtion.Discount
// 			}
// 			if total_quantity > promtion.Quentity {
// 				pot = promotions[len(promotions)-1]
// 				discount_only_product = float32(total_quantity/pot.Quentity) * pot.Discount
// 			}

// 		}

// 		groupedResult = append(groupedResult, custommodel.GroupedByProduct{
// 			ProductID: productID,

// 			PromotionID:         pot.ID,
// 			DiscountOnlyProduct: discount_only_product,
// 			Promotion:           promotions,
// 			TotalPrice:          TotalPrice[productID],
// 			TotalQuantities:     total_quantity,
// 			Items:               groupedItems,
// 		})
// 	}
// 	var orderDetail []*models.OrderDetail
// 	var orderDiscount []*models.OrderDiscount
// 	var orderStockDetails []*models.OrderStockDetail

// 	orderID := uuid.New().String()
// 	var totalPrice float64

// 	for _, data := range groupedResult {

// 		for _, item := range data.Items {
// 			product, err := dbservice.GetProductDetailsForID(db, item.ProductDetailID, item.SizeID)
// 			if err != nil {
// 				return fiber.NewError(http.StatusBadRequest, "not found product")
// 			}
// 			productQuantity := int32(item.Quantity)
// 			for _, stock := range product.StockProducts {

// 				if productQuantity >= 1 {

// 					if stock.Remaining >= productQuantity {
// 						stock.Remaining -= productQuantity

// 						orderStockDetails = append(orderStockDetails, &models.OrderStockDetail{
// 							OrderID:              orderID,
// 							StockProductDetailID: stock.ID,
// 							Quatity:              productQuantity,
// 						})

// 						if stock.Remaining >= 1 {
// 							productQuantity = 0
// 							err = dbservice.UpdateStockProductDetail(db, stock.ID, stock.Remaining, "active")
// 						} else {
// 							err = dbservice.UpdateStockProductDetail(db, stock.ID, stock.Remaining, "out_stock")
// 						}

// 						if err != nil {
// 							return fiber.NewError(http.StatusInternalServerError, "server create order details error")
// 						}

// 					}

// 				}

// 			}
// 			orderDetail = append(orderDetail, &models.OrderDetail{
// 				OrderID:         orderID,
// 				ProductDetailID: item.ProductDetailID,
// 				Quantity:        int32(item.Quantity),
// 				SizeID:          item.SizeID,
// 				UnitPrice:       float64(product.Price),
// 				TotalPrice:      float64((product.Price * item.Quantity)),
// 			})

// 		}

// 		totalPrice += data.TotalPrice

// 		if data.DiscountOnlyProduct > 0 {
// 			orderDiscount = append(orderDiscount, &models.OrderDiscount{
// 				OrderID:       orderID,
// 				ProductID:     data.ProductID,
// 				TotalDiscount: float64(data.DiscountOnlyProduct),
// 				DiscountID:    data.PromotionID,
// 			})
// 		}
// 	}
// 	println(req.FacebookID)

// 	if req.FacebookID != "N/A" && req.PageID != "" {

// 		println(1)

// 		user, _ := dbservice.GetcustomersID(db, req.FacebookID)

// 		// if err != nil {
// 		// 	return fiber.NewError(http.StatusInternalServerError, err.Error())
// 		// }

// 		if user.FacebookID != "" {
// 			println(2)
// 			_, err = dbservice.UpdateColumnsCustomer(db, req.FacebookID, int32(req.Gender), req.Tel)
// 			if err != nil {
// 				return fiber.NewError(http.StatusInternalServerError, err.Error())
// 			}
// 		}

// 		if user.FacebookID == "" {
// 			err := dbservice.CreateCustomer(db, models.Customer{
// 				FacebookID:  req.FacebookID,
// 				FirstName:   req.FullName,
// 				LastName:    "",
// 				Image:       "N/A",
// 				PhoneNumber: req.Tel,
// 				Gender:      int32(req.Gender),
// 				PageID:      req.PageID,
// 			})
// 			if err != nil {
// 				return fiber.NewError(http.StatusInternalServerError, err.Error())
// 			}
// 			println(3)

// 		}

// 	}

// 	// create new customer
// 	if req.FacebookID == "N/A" && req.PageID != "" {

// 		page, err := dbservice.GetPagesByID(db, req.PageID)
// 		req.PlatForm = page.Phalform

// 		if err != nil {
// 			return fiber.NewError(http.StatusInternalServerError, err.Error())
// 		}

// 		if page.Phalform == "facebook" {
// 			return fiber.NewError(400, "ກະລຸນາເລືອກລູກຄ້າຈາກ ລາຍການທີ່ມີຢູ່")
// 		}

// 		customerId := middleware.GenerateFacebookID()
// 		err1 := dbservice.CreateCustomer(db, models.Customer{
// 			FacebookID:  customerId,
// 			FirstName:   req.FullName,
// 			LastName:    "",
// 			Image:       "N/A",
// 			PhoneNumber: req.Tel,
// 			Gender:      int32(req.Gender),
// 			PageID:      req.PageID,
// 		})
// 		if err1 != nil {
// 			return fiber.NewError(http.StatusInternalServerError, err1.Error())
// 		}

// 		req.FacebookID = customerId
// 	}

// 	number := middleware.GenerateOrderNumber()

// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	if req.FacebookID == "" {
// 		return fiber.NewError(http.StatusInternalServerError, " ບໍ່ພົບ id ")

// 	}

// 	order := models.Order{
// 		ID:            orderID,
// 		OrderName:     req.FullName,
// 		Status:        cons.Ordered,
// 		CustomerID:    req.FacebookID,
// 		OrderNo:       number,
// 		Tel:           req.Tel,
// 		ShippingID:    req.ShippingID,
// 		CustomAddress: req.CustomAddress,
// 		UserID:        userID,
// 		DistrictID:    req.DistrictID,
// 		TotalPrice:    totalPrice,
// 		FreeShipping:  req.FreeShipping,
// 		Cod:           req.Cod,
// 		Discount:      req.Discount,
// 		Platform:      req.PlatForm,
// 	}
// 	orderTimeLine := models.OrderTimeLine{
// 		OrderStatus: cons.Ordered,
// 		UserID:      userID,
// 		OrderID:     orderID,
// 	}

// 	err = dbservice.CreateOrder(db, &order)

// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	err = dbservice.CreateOrderTimeLine(db, &orderTimeLine)
// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	err = dbservice.CreateOrderDiscounts(db, orderDiscount, c.Context())
// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	err = dbservice.CreateRoderStockDetails(db, orderStockDetails, c.Context())
// 	if err != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
// 	}

// 	respones, err1 := dbservice.GetOrder(db, orderID)

// 	if err1 != nil {
// 		return fiber.NewError(http.StatusInternalServerError, err1.Error())
// 	}

// 	db.Commit()
// 	webhookURL := "https://discord.com/api/webhooks/1386638914881847366/hz7pb4qbexca75gtnTKNG6G6tLNvlAb5om-21z7ziR_MFvmEkXhKhPLTTPbb4FGtcqH2"
// 	message := fmt.Sprintf(`
// ວັນທີ່: %s
// ລູກຄ້າ: %s
// ລະຫັດ: %s
// ຈັດສົ່ງໂດຍ: %s
// ເບີໂທ: %d
// ທີ່ຢູ່: ເເຂວງ %s ເມືອງ %s ສາຂາ %s
// ຈາກ: %s
// `,
// 		respones.CreatedAt,
// 		respones.OrderName,
// 		respones.OrderNo,
// 		respones.Shipping.Name,
// 		order.Tel,
// 		respones.Province,
// 		respones.District,
// 		order.CustomAddress,
// 		respones.Platform+" "+respones.PageName,
// 	)

// 	SendDiscordWebhook(webhookURL, message)

// 	return c.Status(200).JSON(presenters.ResponseSuccess(respones))
// }

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
	order, err1 := dbservice.GetOrder(db, orderId)
	if err1 != nil {
		return fiber.NewError(500, err1.Error())
	}

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
	webhookURL := "https://discord.com/api/webhooks/1386634660968272027/2MAuHM-iN7ONKqEZ1RjUYQJvdSf51Ck30cb4ojSL5xKY2z7sNXlBLwUqEAqueCui_DTB"
	message := fmt.Sprintf(`
	
ລູກຄ້າ: %s 
ລະຫັດ: %s
ເບີໂທ: %d
ທີ່ຢູ່: ເເຂວງ %s ເມືອງ %s ສາຂາ %s
ຈາກ: %s
`,
		order.OrderName,
		order.OrderNo,
		order.Tel,
		order.Province,
		order.District,
		order.CustomAddress,
		order.Platform+" "+order.PageName,
	)

	SendDiscordWebhook(webhookURL, message)

	return c.Status(200).JSON(presenters.ResponseSuccess("update status success"))
}

func GetSalesHandler(c *fiber.Ctx) error {
	startStr := c.Query("start")
	endStr := c.Query("end")

	// start, err := time.Parse(time.RFC3339, startStr)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"error": "invalid start date"})
	// }
	// end, err := time.Parse(time.RFC3339, endStr)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"error": "invalid end date"})
	// }
	result, err := dbservice.GetProductSalesByHour(gormpkg.GetDB(), startStr, endStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(presenters.ResponseSuccess(result))
}

func CreateOrder(c *fiber.Ctx) error {

	var req custommodel.OrderRequestNew
	var data custommodel.ProductOrderDetails

	var orderProducts []*models.OrderProduct
	var ordProductDetails []*models.OrderProductsDetail
	var ordProductDiscount []*models.OrderProductDiscount
	var orderStockDetails []*models.OrderStockDetail
	var err error
	userID, ok := middleware.GetUserID(c)

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}

	// Parse JSON body into struct
	if err := middleware.ParseAndValidateBody(c, &req); err != nil {
		return fiber.NewError(400, err.Error())
	}

	if !middleware.ValidatePhone(req.Tel) {
		return fiber.NewError(400, "ເບີໂທລະສັບບໍ່ຖືກຕ້ອງ")
	}

	db := gormpkg.GetDB().Begin()
	defer func() {
		db.Rollback()
	}()

	orderID := uuid.New().String()

	for index, prouct := range req.Products {
		orderProductID := uuid.New().String()

		for index1, item := range prouct.ProductDetails {

			productDetail, err1 := dbservice.GetProductDetailsByIDSizdID(db, item.ProductDetailID, item.SizeID, prouct.ProductID)
			if err1 != nil {
				return fiber.NewError(400, err1.Error())
			}
			if productDetail.Remaining < int32(item.Quantity) {
				data, err = dbservice.GetProductByIDSizdID(db, item.ProductDetailID, item.SizeID, prouct.ProductID)
				if err != nil {
					return fiber.NewError(400, err.Error())
				}

				return fiber.NewError(
					550,
					fmt.Sprintf("%s %s %s ຍັງເຫລືອເເຕ່ %v ໂຕ",
						data.Name,
						data.ColorName,
						data.Size,
						productDetail.Remaining,
					),
				)
			}
			req.Products[index].TotalQuantities += item.Quantity
			req.Products[index].TotaProductPrice += float64((item.Quantity * int32(productDetail.Price)))

			// check promotions
			if index1 == (len(prouct.ProductDetails) - 1) {
				productTotalQuantity := req.Products[index].TotalQuantities
				promotions, _ := dbservice.GetPromotion(db, productDetail.ProductID)
				var discount_only_product float32

				for _, promtion := range promotions {

					if productTotalQuantity == promtion.Quentity {
						req.Products[index].Promotion = promtion
						discount_only_product = float32(req.Products[index].TotalQuantities/promtion.Quentity) * promtion.Discount
						if req.Products[index].Promotion.ID != "" {
							ordProductDiscount = append(ordProductDiscount, &models.OrderProductDiscount{
								OrderProductID:   orderProductID,
								PromotionPriceID: req.Products[index].Promotion.ID,
								Discount:         float64(discount_only_product),
							})
						}
					}
					if productTotalQuantity > promtion.Quentity {
						pot := promotions[len(promotions)-1]
						req.Products[index].Promotion = pot
						discount_only_product = float32(req.Products[index].TotalQuantities/promtion.Quentity) * promtion.Discount
						if req.Products[index].Promotion.ID != "" {
							ordProductDiscount = append(ordProductDiscount, &models.OrderProductDiscount{
								OrderProductID:   orderProductID,
								PromotionPriceID: req.Products[index].Promotion.ID,
								Discount:         float64(discount_only_product),
							})
						}
					}

				}
				req.TotalDiscount += prouct.Discount

			}

			ordProductDetails = append(ordProductDetails, &models.OrderProductsDetail{

				ProductDetailID: item.ProductDetailID,
				TotalPrice:      float64((item.Quantity * int32(productDetail.Price))),
				UnitPrice:       productDetail.Price,
				Quantity:        item.Quantity,
				SizeID:          item.SizeID,
				OrderProductID:  orderProductID,
			})
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

		}

		orderProducts = append(orderProducts, &models.OrderProduct{
			OrderID:           orderID,
			ID:                orderProductID,
			Discount:          req.Products[index].Discount,
			ProductID:         prouct.ProductID,
			TotalProductPrice: req.Products[index].TotaProductPrice,
			TotalAmounts:      req.Products[index].TotalQuantities,
		})
		req.TotalPrice += req.Products[index].TotaProductPrice

	}

	// create new customer
	if req.FacebookID != "N/A" && req.PageID != "" {

		println(1)

		user, _ := dbservice.GetcustomersID(db, req.FacebookID)

		if user.FacebookID != "" {

			_, err = dbservice.UpdateColumnsCustomer(db, req.FacebookID, int32(req.Gender), req.Tel)
			if err != nil {
				return fiber.NewError(http.StatusInternalServerError, err.Error())
			}
		}

		if user.FacebookID == "" {
			err := dbservice.CreateCustomer(db, models.Customer{
				FacebookID:  req.FacebookID,
				FirstName:   req.FullName,
				LastName:    "",
				Image:       "N/A",
				PhoneNumber: req.Tel,
				Gender:      int32(req.Gender),
				PageID:      req.PageID,
			})
			if err != nil {
				return fiber.NewError(http.StatusInternalServerError, err.Error())
			}

		}

	}

	// create new customer
	if req.FacebookID == "N/A" && req.PageID != "" {

		page, err := dbservice.GetPagesByID(db, req.PageID)
		req.PlatForm = page.Phalform

		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		if page.Phalform == "facebook" {
			return fiber.NewError(400, "ກະລຸນາເລືອກລູກຄ້າຈາກ ລາຍການທີ່ມີຢູ່")
		}

		customerId := middleware.GenerateFacebookID()
		err1 := dbservice.CreateCustomer(db, models.Customer{
			FacebookID:  customerId,
			FirstName:   req.FullName,
			LastName:    "",
			Image:       "N/A",
			PhoneNumber: req.Tel,
			Gender:      int32(req.Gender),
			PageID:      req.PageID,
		})
		if err1 != nil {
			return fiber.NewError(http.StatusInternalServerError, err1.Error())
		}

		req.FacebookID = customerId
	}

	number := middleware.GenerateOrderNumber()

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
		TotalPrice:    req.TotalPrice,
		FreeShipping:  req.FreeShipping,
		Cod:           req.COD,
		Discount:      req.TotalDiscount,
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

	err = dbservice.CreateOrderProducts(db, orderProducts)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = dbservice.CreateOrderProductDetails(db, ordProductDetails)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if len(ordProductDiscount) > 0 {
		err = dbservice.CreateOrderProductDiscount(db, ordProductDiscount)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	err = dbservice.CreateRoderStockDetails(db, orderStockDetails, c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	err = dbservice.CreateOrderTimeLine(db, &orderTimeLine)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	orderReponse, err1 := dbservice.GetOrder(db, orderID)
	if err1 != nil {
		return fiber.NewError(500, err.Error())
	}
	db.Commit()
	// 	webhookURL := "https://discord.com/api/webhooks/1386638914881847366/hz7pb4qbexca75gtnTKNG6G6tLNvlAb5om-21z7ziR_MFvmEkXhKhPLTTPbb4FGtcqH2"
	// 	message := fmt.Sprintf(`
	// ວັນທີ່: %s
	// ລູກຄ້າ: %s
	// ລະຫັດ: %s
	// ຈັດສົ່ງໂດຍ: %s
	// ເບີໂທ: %d
	// ທີ່ຢູ່: ເເຂວງ %s ເມືອງ %s ສາຂາ %s
	// ຈາກ: %s
	// `,
	// 		respones.CreatedAt,
	// 		respones.OrderName,
	// 		respones.OrderNo,
	// 		respones.Shipping.Name,
	// 		order.Tel,
	// 		respones.Province,
	// 		respones.District,
	// 		order.CustomAddress,
	// 		respones.Platform+" "+respones.PageName,
	// 	)

	// 	SendDiscordWebhook(webhookURL, message)

	return c.Status(200).JSON(presenters.ResponseSuccess(orderReponse))
}
