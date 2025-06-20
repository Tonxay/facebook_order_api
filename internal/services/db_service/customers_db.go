package dbservice

import (
	"encoding/json"
	"fmt"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"
	"io"
	"net/http"
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

func GetUserInFaceBook(pageId, pageAccessToken string) (custommodel.FacebookConversationList, error) {
	var result custommodel.FacebookConversationList

	if pageId == "" || pageAccessToken == "" {
		return result, fmt.Errorf("pageId or pageAccessToken is missing")
	}

	// Build URL (no token in URL now)
	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/%s/conversations?fields=participants,updated_time&limit=50",
		pageId,
	)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %v", err)
	}

	// Set Authorization header
	req.Header.Set("Authorization", "Bearer "+pageAccessToken)

	// Send request using default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("non-200 response from Facebook API: %d, %s", resp.StatusCode, string(body))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to parse JSON response: %v\nraw: %s", err, string(body))
	}

	return result, nil
}
