package dbservice

import (
	"context"
	"fmt"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/query"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *models.Order) error {
	err := db.Table(models.TableNameOrder).Create(order).Error
	return err
}

func CreateOrderDetails(db *gorm.DB, orderDetail []*models.OrderDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.OrderDetail

	// Insert in batches, handle error
	err := daq.WithContext(ctx).CreateInBatches(orderDetail, 100)
	if err != nil {
		// Optional: Log or wrap for context
		return fmt.Errorf("failed to create order details: %w", err)
	}

	return nil
}

func CreateOrderDiscounts(db *gorm.DB, orderDiscounts []*models.OrderDiscount, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.OrderDiscount
	// Insert in batches, handle error
	err := daq.WithContext(ctx).CreateInBatches(orderDiscounts, 100)
	if err != nil {
		// Optional: Log or wrap for context
		return fmt.Errorf("failed to create order details: %w", err)
	}
	return nil
}

func GetOrders(db *gorm.DB, statuses []string, isCancell bool) ([]*custommodel.OrderReponse, error) {
	var orders []*custommodel.OrderReponse

	tx := db.Table(models.TableNameOrder + " o").Select(`o.*,SUM(rc.total_discount) AS total_prodouct_discount,d.dr_name,provice.pr_name,page.name_page AS page_name`)

	if !isCancell {
		tx = tx.Where("status IN ?", statuses)
	}

	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderDiscount + " rc ON rc.order_id = o.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")

	tx = tx.Where("o.is_cancel = ?", isCancell)

	tx = tx.Preload("OrderDetails").Preload("OrderDetails.ProductDetail").Preload("OrderDetails.ProductDetail.Product").
		Preload("OrderDetails.Size").Preload("Shipping")

	tx = tx.Preload("OrderDiscounts")

	tx = tx.Group(`o.id,d.dr_name,provice.pr_name,page.name_page`)

	err := tx.Order("o.updated_at DESC").Find(&orders).Error
	return orders, err
}
func GetOrder(db *gorm.DB, orderID string) (custommodel.OrderReponse, error) {
	var orders custommodel.OrderReponse

	tx := db.Table(models.TableNameOrder+" o").Select(`o.*,SUM(rc.total_discount) AS total_prodouct_discount,d.dr_name,provice.pr_name, page.name_page AS page_name`).Where("o.id = ?", orderID)

	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderDiscount + " rc ON rc.order_id = o.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")
	// tx = tx.Joins("LEFT JOIN " + models.TableNameOrderTimeLine + " orl ON orl.order_id = o.id")

	// tx = tx.Where("orl.order_status != ?", cons.OrderCancelled)

	tx = tx.Preload("OrderDetails").Preload("OrderDetails.ProductDetail").Preload("OrderDetails.ProductDetail.Product").
		Preload("OrderDetails.Size").Preload("Shipping")

	tx = tx.Preload("OrderDiscounts")

	tx = tx.Group(`o.id,d.dr_name,provice.pr_name,page.name_page`)

	err := tx.Find(&orders).Error
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

func UpdateStatusOrder(db *gorm.DB, orderId, newStatus, oldStatus string, orderNo string) (models.Order, error) {
	var order models.Order
	tx := db.Table(models.TableNameOrder)
	if orderId != "" {
		tx = tx.Where("id = ? AND status = ? AND is_cancel = ?", orderId, oldStatus, false)
	} else {
		tx = tx.Where("order_no = ? AND status = ? AND is_cancel = ?", orderNo, oldStatus, false)
	}

	result := tx.UpdateColumns(&models.Order{
		Status: newStatus,
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

func UpdateIsCancelOrder(db *gorm.DB, orderId string) error {

	result := db.Table(models.TableNameOrder).Where("id = ? AND is_cancel = ?", orderId, false).UpdateColumns(&models.Order{
		IsCancel: true,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated — order may not exist or status is unchanged")
	}

	return result.Error
}
