package dbservice

import (
	gormpkg "github.com/yourusername/go-api/internal/pkg"
	"github.com/yourusername/go-api/internal/pkg/models"
)

func Getcustomers() (*[]models.Customer, error) {
	var customers *[]models.Customer
	err := gormpkg.GetDB().Table(models.TableNameCustomer).Find(&customers).Error
	return customers, err
}
