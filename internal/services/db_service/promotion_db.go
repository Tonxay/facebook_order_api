package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func GetPromotion(db *gorm.DB, productId string) ([]models.Promotion, error) {
	promotion := []models.Promotion{}
	err := db.Table(models.TableNamePromotion).Where("status = ? AND product_id = ?", "active", productId).Find(&promotion).Error
	return promotion, err
}
