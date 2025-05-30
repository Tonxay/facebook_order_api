package service

import (
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	dbservice "go-api/internal/service/db_service"

	"github.com/gofiber/fiber/v2"
)

func GetFacebookProfile(c *fiber.Ctx) error {
	id := c.Params("facebook_id")
	user, err := getFacebookProfile(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "create messeng",
		})
	}

	return c.JSON(user)
}
func GetFacebookAllCustomers(c *fiber.Ctx) error {

	user, err := dbservice.Getcustomers(gormpkg.GetDB())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "request error",
		})
	}

	return c.JSON(user)
}

func getFacebookProfile(facebookID string) (*models.Customer, error) {

	var customer models.Customer
	gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", facebookID).First(&customer)

	return &customer, nil
}
