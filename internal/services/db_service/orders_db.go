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

func GetOrders(db *gorm.DB) ([]*custommodel.OrderReponse, error) {
	var orders []*custommodel.OrderReponse
	tx := db.Table(models.TableNameOrder + " o").Select(`o.*,SUM(rc.total_discount) AS total_prodouct_discount,d.dr_name,provice.pr_name,page.name_page AS page_name`)

	tx = tx.Joins("LEFT JOIN " + models.TableNameOrderDiscount + " rc ON rc.order_id = o.id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameDistrict + " d ON d.id = o.district_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameProvince + " provice ON provice.id = d.province_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNameCustomer + " c ON c.facebook_id = o.customer_id")
	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = c.page_id")

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

func UpdateOrder(db *gorm.DB, orderId, newStatus, oldStatus string) error {
	result := db.Table(models.TableNameOrder).Where("id = ? AND status = ?", orderId, oldStatus).UpdateColumns(&models.Order{
		Status: newStatus,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated — order may not exist or status is unchanged")
	}

	return result.Error
}
