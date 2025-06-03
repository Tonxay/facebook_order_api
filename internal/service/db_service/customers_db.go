package dbservice

import (
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/request"

	"gorm.io/gorm"
)

func Getcustomers(db *gorm.DB, query request.CustomerQuery) (*[]models.Customer, int64, error) {
	var totalCount int64
	var customers *[]models.Customer
	offset := (query.Page - 1) * query.Limit
	err := db.
		Table(models.TableNameCustomer).
		Count(&totalCount).
		Limit(query.Limit).
		Offset(offset).
		Find(&customers).Error
	return customers, totalCount, err
}
