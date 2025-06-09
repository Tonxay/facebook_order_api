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

func GetProvince(db *gorm.DB) (*[]models.Province, error) {
	var provices *[]models.Province
	err := db.Table(models.TableNameProvince).Find(&provices).Error
	return provices, err
}

func GetDistrict(db *gorm.DB, id string) (*[]models.District, error) {
	var dis *[]models.District
	err := db.Table(models.TableNameDistrict).Where("province_id = ?", id).Find(&dis).Error
	return dis, err
}
