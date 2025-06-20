package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func GetCategory(db *gorm.DB) ([]models.Category, error) {
	var cates []models.Category
	err := db.Model(&cates).Where("status = ?", "active").Find(&cates).Error
	return cates, err
}
