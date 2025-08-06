package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func CreateSize(db *gorm.DB, size models.Size) (models.Size, error) {
	err := db.Table(models.TableNameSize).Create(&size).Where("product_detail_id = ?", size.ProductDetailID).First(&size).Error
	return size, err
}
