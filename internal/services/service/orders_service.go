package service

import (
	"fmt"
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/services/db_service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {

	var req custommodel.OrderRequest

	// var err error

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
	total_quantity := make(map[string]int)
	TotalPrice := make(map[string]int32)
	price := make(map[string]int)

	for _, item := range req.Items {

		data, err := dbservice.GetProductDetailsByIDSizdID(db, item.ProductDetailID, item.SizeID, item.ProductID)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		if data.Remaining < int32(item.Quantity) {

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
		TotalPrice[item.ProductID] += (data.Price * int32(item.Quantity))
	}

	var discount_only_product int
	// Step 2: Build final output
	var groupedResult []custommodel.GroupedByProduct
	for productID, groupedItems := range temp {
		promotions, _ := dbservice.GetPromotion(db, productID)
		total_quantity := total_quantity[productID]

		for _, promtion := range promotions {
			if total_quantity == int(promtion.Quentity) {
				discount_only_product = (total_quantity / int(promtion.Quentity)) * int(promtion.Discount)
			}
			if total_quantity > int(promtion.Quentity) {
				pot := promotions[len(promotions)-1]
				discount_only_product = (total_quantity / int(pot.Quentity)) * int(pot.Discount)
			}

		}

		groupedResult = append(groupedResult, custommodel.GroupedByProduct{
			ProductID:           productID,
			DiscountOnlyProduct: discount_only_product,
			Promotion:           promotions,
			TotalPrice:          TotalPrice[productID],
			TotalQuantities:     total_quantity,
			Items:               groupedItems,
		})
	}

	// Step 3: Marshal and print
	// result, err := json.MarshalIndent(groupedResult, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }

	// var orderDetail []*models.OrderDetail

	// orderID := uuid.New().String()
	// var totalPrice int32

	// for _, item := range req.OrderDetails {
	// 	product, err := dbservice.GetProductDetailsForID(db, item.ProductDetailID, item.SizeID)
	// 	if err != nil {
	// 		return fiber.NewError(http.StatusBadRequest, "not found product")

	// 	}

	// 	if item.Quantity > product.Quantity {
	// 		return fiber.NewError(http.StatusBadRequest, "quantity than product")

	// 	}

	// productQuantity := item.Quantity
	// for _, stock := range product.StockProducts {

	// 	if stock.Remaining <= 0 {
	// 		continue
	// 	}
	// 	if productQuantity >= 1 {

	// 		if stock.Remaining >= productQuantity {
	// 			stock.Remaining -= productQuantity
	// 			productQuantity = 0
	// 			err = dbservice.UpdateStockProductDetail(stock.ID, stock.Remaining, "active", c.Context())
	// 			if err != nil {
	// 				return fiber.NewError(http.StatusInternalServerError, "server create order details error")
	// 			}
	// 		} else {
	// 			productQuantity -= stock.Remaining
	// 			stock.Remaining = 0
	// 			err = dbservice.UpdateStockProductDetail(stock.ID, stock.Remaining, "out_stock", c.Context())
	// 			if err != nil {
	// 				return fiber.NewError(http.StatusInternalServerError, "server create order details error")

	// 			}
	// 		}

	// 	}

	// }

	// 	orderDetail = append(orderDetail, &models.OrderDetail{
	// 		OrderID:         orderID,
	// 		ProductDetailID: item.ProductDetailID,
	// 		Quantity:        item.Quantity,
	// 		SizeID:          item.SizeID,
	// 		UnitPrice:       float64(product.Price),
	// 		TotalPrice:      float64(item.Quantity * product.Price),
	// 	})

	// 	totalPrice = (item.Quantity * product.Price) + totalPrice

	// }

	// order := models.Order{
	// 	ID:            orderID,
	// 	Status:        "pending",
	// 	CustomerID:    req.CustomerID,
	// 	PackagePrice:  0,
	// 	OrderNo:       middleware.GenerateOrderNumber(),
	// 	Tel:           req.Tel,
	// 	CustomAddress: req.CustomAddress,
	// 	UserID:        "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77",
	// 	DistrictID:    req.DistrictID,
	// 	TotalPrice:    totalPrice,
	// }

	// err = dbservice.CreateOrder(db, &order, c.Context())
	// if err != nil {

	// 	return fiber.NewError(http.StatusInternalServerError, "server create order details error")

	// }

	// err = dbservice.CreateOrderDetails(db, orderDetail, c.Context())
	// if err != nil {

	// 	return fiber.NewError(http.StatusInternalServerError, "server create order details error")
	// }
	db.Commit()

	return c.Status(200).JSON(presenters.ResponseSuccess(fiber.Map{
		"order": groupedResult,
		// "items": orderDetail,
	}))

}
