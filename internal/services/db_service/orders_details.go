package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

// func GetOrderDetails(db *gorm.DB, orderID string) ([]*models.OrderDetail, error) {
// 	var orderDetails []*models.OrderDetail

// 	tx := db.Table(models.TableNameOrderDetail + " AS pd")
// 	tx = tx.Joins("LEFT JOIN " + models.TableNameOrder + "  ord ON ord.id = pd.order_id")
// 	tx = tx.Where("pd.order_id = ? ", orderID)
// 	tx = tx.Where("ord.is_cancel = ?", false)
// 	err := tx.Find(&orderDetails).Error

// 	return orderDetails, err

// }

func GetOrderDetails(db *gorm.DB, orderID string) ([]models.OrderProductsDetail, error) {
	var orderDetails []models.OrderProductsDetail

	tx := db.Table(models.TableNameOrderProductsDetail+" AS orpd").
		Select("orpd.*"). // ✅ Select fields matching the destination struct
		Joins("LEFT JOIN "+models.TableNameOrderProduct+" AS orp ON orp.id = orpd.order_product_id").
		Joins("LEFT JOIN "+models.TableNameOrder+" AS ord ON ord.id = orp.order_id").
		Where("ord.id = ?", orderID).
		Where("ord.is_cancel = ?", false)

	err := tx.Find(&orderDetails).Error
	return orderDetails, err
}
