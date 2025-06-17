package dbservice

import (
	"go-api/internal/pkg/models"

	"gorm.io/gorm"
)

func GetPages(db *gorm.DB) ([]models.Page, error) {
	var pages []models.Page
	tx := db.Table(models.TableNamePage)
	err := tx.Where("status = ?", 1).Order("name_page ASC").Find(&pages).Error
	return pages, err
}

func GetPagesByID(db *gorm.DB, pageID string) (models.Page, error) {
	var pages models.Page
	tx := db.Table(models.TableNamePage)
	err := tx.Where("status = ? AND page_id = ?", 1, pageID).Order("name_page ASC").Find(&pages).Error
	return pages, err
}
