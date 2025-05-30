package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func GetShipping(db *gorm.DB) (*[]models.Shipping, error) {
	var ship *[]models.Shipping
	err := db.Table(models.TableNameShipping).Where("status = ?", "active").Find(&ship).Error
	return ship, err
}
