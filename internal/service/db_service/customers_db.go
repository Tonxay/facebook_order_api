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
	tx := db.
		Table(models.TableNameCustomer)
	tx = tx.Count(&totalCount)
	if query.Name != "" {
		tx = tx.Where("first_name ILIKE ?", "%"+query.Name+"%")
	}

	tx = tx.Limit(query.Limit).
		Offset(offset)

	err := tx.Order("created_at DESC").Find(&customers).Error
	return customers, totalCount, err

}
