package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/services/db_service"

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
func GetUserInFaceBookJson(c *fiber.Ctx) error {
	id := c.Query("page_id")

	pageId, token := middleware.CheckPageId(id, id)

	user, err := dbservice.GetUserInFaceBook(pageId, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "create messeng",
		})
	}

	return c.JSON(presenters.ResponseSuccess(user))
}

func GetFacebookAllCustomers(c *fiber.Ctx) error {
	query := request.CustomerQuery{}

	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(400, "request erro")
	}
	query.Limit = c.QueryInt("limit", 20)
	query.Page = c.QueryInt("page", 1)

	user, total, err := dbservice.Getcustomers(gormpkg.GetDB(), query)
	if err != nil {
		return fiber.NewError(500, "request erro")
	}

	totalPage := int((total + int64(query.Limit) - 1) / int64(query.Limit))

	return c.JSON(presenters.ResponseSuccessListData(user, query.Page, query.Limit, int(total), totalPage))
}

func getFacebookProfile(facebookID string) (*models.Customer, error) {

	var customer models.Customer
	gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", facebookID).First(&customer)

	return &customer, nil
}
