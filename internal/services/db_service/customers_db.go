package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"

	"gorm.io/gorm"
)

func Getcustomers(db *gorm.DB, query request.CustomerQuery) (*[]custommodel.Customer, int64, error) {
	var totalCount int64
	var customers *[]custommodel.Customer
	offset := (query.Page - 1) * query.Limit
	tx := db.
		Table(models.TableNameCustomer + " cus")

	tx = tx.Select(` cus.* ,page.* , page.image AS page_image `)

	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = cus.page_id")

	tx = tx.Count(&totalCount)
	if query.Name != "" {
		tx = tx.Where("cus.first_name ILIKE ?", "%"+query.Name+"%")
	} else {
		tx = tx.Limit(query.Limit).
			Offset(offset)
	}

	err := tx.Order("updated_at DESC").Find(&customers).Error
	return customers, totalCount, err

}
