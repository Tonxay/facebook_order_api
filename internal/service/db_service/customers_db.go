package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func Getcustomers(db *gorm.DB) (*[]models.Customer, error) {
	var customers *[]models.Customer
	err := db.Table(models.TableNameCustomer).Find(&customers).Error
	return customers, err
}
