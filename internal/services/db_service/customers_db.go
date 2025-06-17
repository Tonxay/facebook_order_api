package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"
	"time"

	"gorm.io/gorm"
)

func Getcustomers(db *gorm.DB, query request.CustomerQuery) (*[]custommodel.Customer, int64, error) {
	var totalCount int64
	var customers *[]custommodel.Customer
	offset := (query.Page - 1) * query.Limit
	tx := db.
		Table(models.TableNameCustomer + " cus")

	tx = tx.Select(` cus.* ,page.* `)

	tx = tx.Joins("LEFT JOIN " + models.TableNamePage + " page ON page.page_id = cus.page_id")

	tx = tx.Count(&totalCount)
	tx = tx.Preload("Page")
	if query.Name != "" {
		tx = tx.Where("cus.first_name ILIKE ?", "%"+query.Name+"%")
	} else {
		tx = tx.Limit(query.Limit).
			Offset(offset)
	}

	err := tx.Order("updated_at DESC").Find(&customers).Error
	return customers, totalCount, err

}
func GetcustomersID(db *gorm.DB, fbID string) (custommodel.Customer, error) {
	var user custommodel.Customer
	err := db.Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).First(&user).Error
	return user, err
}
func UpdateColumnsCustomer(db *gorm.DB, fbID string, gender int32, tel int64) (models.Customer, error) {
	var user models.Customer
	err := db.Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).UpdateColumns(&models.Customer{
		Gender:      gender,
		PhoneNumber: tel,
		UpdatedAt:   time.Now(),
	}).First(&user).Error
	return user, err
}
func CreateColumnsCustomer(db *gorm.DB, fbID string, pageId string) (custommodel.Customer, error) {
	var user custommodel.Customer
	err := db.Table(models.TableNameCustomer).Create(&models.Customer{
		FacebookID: fbID,
		PageID:     pageId,
	}).First(&user).Error
	return user, err
}

func CreateCustomer(db *gorm.DB, newCustomer models.Customer) (models.Customer, error) {
	var user models.Customer
	err := db.Table(models.TableNameCustomer).Create(&newCustomer).First(&user).Error
	return user, err
}
