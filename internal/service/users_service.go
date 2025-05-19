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
	// pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	// if pageAccessToken == "" {
	// 	return nil, fmt.Errorf("missing PAGE_ACCESS_TOKEN")
	// }

	// url := fmt.Sprintf(
	// 	"https://graph.facebook.com/v21.0/%s?fields=first_name,last_name,name,email,gender,locale,timezone,profile_pic&access_token=%s",
	// 	facebookID,
	// 	pageAccessToken,
	// )

	// resp, err := http.Get(url)
	// if err != nil {
	// 	return nil, fmt.Errorf("error making request to Facebook Graph API: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("facebook API error: %s", resp.Status)
	// }

	// var profile custommodel.FacebookUserProfile
	// if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
	// 	return nil, fmt.Errorf("error decoding Facebook response: %w", err)
	// }
	var customer models.Customer
	gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", facebookID).First(&customer)

	return &customer, nil
}
