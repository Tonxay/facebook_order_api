package dbservice

import (
	"fmt"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"

	"gorm.io/gorm"
)

func GetProductDetailsForID(db *gorm.DB, productItemDetailId string) (custommodel.ProductDetailCounter, error) {
	var products custommodel.ProductDetailCounter

	tx := db.Table(models.TableNameProductDetail + " pd")

	tx = tx.Select("pd.id,p.name, pd.size, pd.color, SUM(spd.quantity) AS quantity , pd.price")

	tx = tx.Where("pd.id = ?", productItemDetailId)

	joinProduct := fmt.Sprintf("LEFT JOIN %s p ON pd.product_id = p.id", models.TableNameProduct)
	joinQuery := fmt.Sprintf("LEFT JOIN %s spd ON spd.product_detail_id = pd.id", models.TableNameStockProductDetail)

	tx = tx.Joins(joinQuery).Joins(joinProduct)

	tx = tx.Preload("StockProducts", func(db *gorm.DB) *gorm.DB {
		return db.Where("remaining > ? AND status = ?", 0, "active")
	})

	tx = tx.Group("pd.id,p.name, pd.size, pd.color,pd.price")

	err := tx.Find(&products).Error
	return products, err
}
