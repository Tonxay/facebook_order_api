package dbservice

import (
	"context"
	"fmt"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"
	"go-api/internal/pkg/query"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *models.Order) error {
	err := db.Table(models.TableNameOrder).Create(order).Error
	return err
}

// func CreateOrderDetails(db *gorm.DB, orderDetail []*models.OrderDetail, ctx context.Context) error {
// 	query.SetDefault(db)
// 	daq := query.Q.OrderDetail

// 	// Insert in batches, handle error
// 	err := daq.WithContext(ctx).CreateInBatches(orderDetail, 100)
// 	if err != nil {
// 		// Optional: Log or wrap for context
// 		return fmt.Errorf("failed to create order details: %w", err)
// 	}

// 	return nil
// }

func CreateOrderDiscounts(db *gorm.DB, orderDiscounts []*models.OrderDiscount, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.OrderDiscount
	// Insert in batches, handle error

	err := daq.WithContext(ctx).Where(daq.DiscountID).CreateInBatches(orderDiscounts, 100)
	if err != nil {
		// Optional: Log or wrap for context
		return fmt.Errorf("failed to create order details: %w", err)
	}
	return nil
}

// func GetOrders(db *gorm.DB, filter request.StatusOrderRequest) ([]*custommodel.OrderReponse, error) {
// 	var orders []*custommodel.OrderReponse

// 	tx := db.Table(models.TableNameOrder + " o").Select(
// 		`o.*,
// 	SUM(rc.total_discount) AS total_prodouct_discount,
// 	d.dr_name,provice.pr_name,
// 	page.name_page AS page_name,
// 	page.tel AS page_tel`)

// 	if !filter.IsCancel {
// 		tx = tx.Where("o.status IN ?", filter.Statuses)
// 	}

// 	if filter.Tel != "" {
// 		tx = tx.Where("o.tel::text LIKE ?", "%"+filter.Tel+"%")
// 	}

// 	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderDiscount + " rc ON rc.order_id = o.id")
// 	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
// 	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
// 	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
// 	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")

// 	tx = tx.Where("o.is_cancel = ?", filter.IsCancel)

// 	tx = tx.Preload("OrderDetails").Preload("OrderDetails.ProductDetail").Preload("OrderDetails.ProductDetail.Product").
// 		Preload("OrderDetails.Size").Preload("Shipping")

// 	tx = tx.Preload("OrderDiscounts")

// 	tx = tx.Group(`o.id,d.dr_name,provice.pr_name,page.name_page,page.tel`)

// 	err := tx.Order("o.updated_at DESC").Find(&orders).Error
// 	return orders, err
// }

func GetOrders(db *gorm.DB, filter request.StatusOrderRequest) ([]*custommodel.OrderReponseNew, error) {
	var orders []*custommodel.OrderReponseNew

	tx := db.Table(models.TableNameOrder + " o").Select(
		`o.*,
	d.dr_name,provice.pr_name,
	SUM(rc.discount) AS total_prodouct_discount,
	page.name_page AS page_name,
	page.tel AS page_tel`)

	if !filter.IsCancel {
		tx = tx.Where("o.status IN ?", filter.Statuses)
	}

	if filter.Tel != "" {
		tx = tx.Where("o.tel::text LIKE ?", "%"+filter.Tel+"%")
	}

	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderProduct + " op ON op.order_id = o.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderProductDiscount + " rc ON rc.order_product_id = op.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")

	tx = tx.Where("o.is_cancel = ?", filter.IsCancel)

	tx = tx.Preload("OrderProducts").Preload("OrderProducts.Product").Preload("OrderProducts.OrderProductsDetails").Preload("OrderProducts.OrderProductDiscount").
		Preload("OrderProducts.OrderProductsDetails.Size").Preload("OrderProducts.OrderProductsDetails.ProductDetail").Preload("Shipping")

	// tx = tx.Preload("OrderDiscounts")

	tx = tx.Group(`o.id,d.dr_name,provice.pr_name,page.name_page,page.tel`)

	err := tx.Order("o.updated_at DESC").Find(&orders).Error
	return orders, err
}
func GetOrder(db *gorm.DB, orderID string) (custommodel.OrderReponseNew, error) {
	var orders custommodel.OrderReponseNew

	tx := db.Table(models.TableNameOrder+" o").Select(
		`o.*,
	d.dr_name,provice.pr_name,
	SUM(rc.discount) AS total_prodouct_discount,
	page.name_page AS page_name,
	page.tel AS page_tel
	`,
	).Where("o.id = ?", orderID)

	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderProduct + " op ON op.order_id = o.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderProductDiscount + " rc ON rc.order_product_id = op.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")

	tx = tx.Preload("OrderProducts").Preload("OrderProducts.Product").Preload("OrderProducts.OrderProductsDetails").Preload("OrderProducts.OrderProductDiscount").
		Preload("OrderProducts.OrderProductsDetails.Size").Preload("OrderProducts.OrderProductsDetails.ProductDetail").Preload("Shipping")

	tx = tx.Group(`o.id,d.dr_name,provice.pr_name,page.name_page,page.tel`)

	err := tx.Order("o.updated_at DESC").Find(&orders).Error
	return orders, err
}

func CreateOrderTimeLine(db *gorm.DB, orderTimeLine *models.OrderTimeLine) error {
	result := db.Table(models.TableNameOrderTimeLine).Create(orderTimeLine)
	if result.Error != nil {
		return fmt.Errorf("failed to create order time line: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to create order time line")
	}
	return result.Error
}

func UpdateStatusOrder(db *gorm.DB, orderId, newStatus, oldStatus string, orderNo string, userId string) (models.Order, error) {
	var order models.Order
	tx := db.Table(models.TableNameOrder)
	if orderId != "" {
		tx = tx.Where("id = ? AND status = ? AND is_cancel = ?", orderId, oldStatus, false)
	} else {
		tx = tx.Where("order_no = ? AND status = ? AND is_cancel = ?", orderNo, oldStatus, false)
	}

	result := tx.UpdateColumns(&models.Order{
		Status:      newStatus,
		UserUpdated: userId,
	})

	if result.Error != nil {
		return order, fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return order, fmt.Errorf("no rows updated — order may not exist or status is unchanged")
	}
	tx1 := db.Table(models.TableNameOrder)

	if orderId != "" {
		tx1 = tx1.Where("id = ? AND status = ? AND is_cancel = ?", orderId, newStatus, false)
	} else {
		tx1 = tx1.Where("order_no = ? AND status = ? AND is_cancel = ?", orderNo, newStatus, false)
	}

	err := tx1.First(&order).Error
	if err != nil {
		return order, err
	}
	return order, result.Error
}

func UpdateIsCancelOrder(db *gorm.DB, orderId string, userID string) error {

	result := db.Table(models.TableNameOrder).Where("id = ? AND is_cancel = ?", orderId, false).UpdateColumns(&models.Order{
		IsCancel:    true,
		UserUpdated: userID,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated — order may not exist or status is unchanged")
	}

	return result.Error
}

func GetProductSalesByHour(db *gorm.DB, startDate string, endDate string) ([]custommodel.ProductSalesGrouped, error) {

	var flatData []custommodel.FlatProductSales

	tx := db.
		Model(&custommodel.OrderDetail{}).
		Select(`
			p.name AS product_name,
			to_char(d.created_at, 'HH24') AS hour,
			SUM(order_details.quantity) AS total_quantity,
			SUM(order_details.total_price) AS total_price
		`).
		Joins("JOIN product_details pd ON pd.id = order_details.product_detail_id").
		Joins("JOIN products p ON pd.product_id = p.id").
		Joins("JOIN orders d ON order_details.order_id = d.id").
		Joins("LEFT JOIN "+models.TableNameOrderDiscount+" dc ON dc.order_id = d.id").
		Where("d.is_cancel = ?", false)

	if startDate != "" && endDate != "" {
		tx = tx.Where("d.ordered_at BETWEEN ? AND ?", startDate, endDate)
	}

	tx = tx.Group("p.name, to_char(d.created_at, 'HH24')")
	tx = tx.Order("p.name, hour")

	err := tx.Scan(&flatData).Error

	if err != nil {
		return nil, err
	}

	// Group by product name
	groupMap := make(map[string][]custommodel.HourlyStat)

	for _, record := range flatData {
		groupMap[record.ProductName] = append(groupMap[record.ProductName], custommodel.HourlyStat{
			Hour:       record.Hour,
			Quantity:   record.Quantity,
			TotalPrice: record.TotalPrice,
		})
	}

	// Convert map to slice
	var result []custommodel.ProductSalesGrouped
	for name, times := range groupMap {
		result = append(result, custommodel.ProductSalesGrouped{
			ProductName: name,
			Times:       times,
		})
	}

	return result, nil
}

func CreateOrderProducts(db *gorm.DB, orderProduct []*models.OrderProduct) error {
	result := db.Table(models.TableNameOrderProduct).CreateInBatches(orderProduct, 50)

	if result.Error != nil {
		return fmt.Errorf("failed to create order product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to create order product")
	}

	return result.Error
}

func CreateOrderProductDetails(db *gorm.DB, productsDetails []*models.OrderProductsDetail) error {
	// Perform batch insert
	result := db.Table(models.TableNameOrderProductsDetail).CreateInBatches(productsDetails, 50)

	if result.Error != nil {
		return fmt.Errorf("failed to create order product details: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to create order product details")
	}

	return nil
}

func CreateOrderProductDiscount(db *gorm.DB, productsDiscount []*models.OrderProductDiscount) error {
	result := db.Table(models.TableNameOrderProductDiscount).CreateInBatches(productsDiscount, 20)

	if result.Error != nil {
		return fmt.Errorf("failed to create order product discount: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to create order product discount")
	}

	return result.Error
}
