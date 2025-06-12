package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func GetOrderDetails(db *gorm.DB, orderID string) ([]*models.OrderDetail, error) {
	var orderDetails []*models.OrderDetail

	tx := db.Table(models.TableNameOrderDetail + " AS pd")
	tx = tx.Joins("LEFT JOIN " + models.TableNameOrder + "  ord ON ord.id = pd.order_id")
	tx = tx.Where("pd.order_id = ? ", orderID)
	tx = tx.Where("ord.is_cancel = ?", false)
	err := tx.Find(&orderDetails).Error

	return orderDetails, err

}
