package dbservice

import (
	"fmt"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"

	"gorm.io/gorm"
)

func GetProductDetailsForID(db *gorm.DB, productItemDetailId string, size_id string) (custommodel.ProductDetailCounter, error) {
	var products custommodel.ProductDetailCounter

	tx := db.Table(models.TableNameProductDetail + " pd")

	tx = tx.Select("pd.id,p.name, s.size, pd.color, SUM(spd.remaining) AS quantity ")

	tx = tx.Where("pd.id = ?", productItemDetailId).Where("spd.status = ?", "active").Where("size_id = ?", size_id)

	joinProduct := fmt.Sprintf("LEFT JOIN %s p ON pd.product_id = p.id", models.TableNameProduct)
	joinQuery := fmt.Sprintf("LEFT JOIN %s spd ON spd.product_detail_id = pd.id", models.TableNameStockProductDetail)
	joinSize := fmt.Sprintf("LEFT JOIN %s s ON s.product_detail_id = pd.id", models.TableNameSize)

	tx = tx.Joins(joinQuery).Joins(joinProduct).Joins(joinSize)

	tx = tx.Preload("StockProducts", func(db *gorm.DB) *gorm.DB {
		return db.Where("remaining > ? AND status = ? AND size_id = ?", 0, "active", size_id)
	})

	tx = tx.Group("pd.id,p.name, s.size, pd.color")

	err := tx.Find(&products).Error

	return products, err
}

func GetProductDetails(db *gorm.DB, productItemDetailId string) (custommodel.ProductDetails, error) {
	var products custommodel.ProductDetails

	tx := db.Table(models.TableNameProductDetail + " pd")

	tx = tx.Select("pd.id,p.name, pd.size, pd.color, SUM(spd.remaining) AS quantity , pd.price")

	tx = tx.Where("pd.id = ?", productItemDetailId).Where("spd.status = ?", "active")

	joinProduct := fmt.Sprintf("LEFT JOIN %s p ON pd.product_id = p.id", models.TableNameProduct)
	joinQuery := fmt.Sprintf("LEFT JOIN %s spd ON spd.product_detail_id = pd.id", models.TableNameStockProductDetail)

	tx = tx.Joins(joinQuery).Joins(joinProduct)

	// tx = tx.Preload("StockProducts", func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("remaining > ? AND status = ?", 0, "active")
	// })

	tx = tx.Group("pd.id,p.name, pd.size, pd.color,pd.price")

	err := tx.Find(&products).Error
	return products, err
}

func GetProductDetailsByIDSizdID(db *gorm.DB, productItemDetailId, sizeId, producId string) (custommodel.ProductOrderDetails, error) {
	var products custommodel.ProductOrderDetails

	tx := db.Table(models.TableNameProductDetail + " pd")

	tx = tx.Select("pd.product_id, pd.id, p.name, s.id  AS size_id , s.size ,pd.color_name, SUM(spd.remaining) AS remaining , s.price")

	tx = tx.Where("pd.id = ?", productItemDetailId).Where("spd.status = ?", "active")

	joinProduct := fmt.Sprintf("LEFT JOIN %s p ON pd.product_id = p.id", models.TableNameProduct)
	joinQuery := fmt.Sprintf("LEFT JOIN %s spd ON spd.product_detail_id = pd.id", models.TableNameStockProductDetail)
	joinSize := fmt.Sprintf("LEFT JOIN %s s ON s.product_detail_id = pd.id", models.TableNameSize)

	tx = tx.Joins(joinQuery).Joins(joinProduct).Joins(joinSize)

	tx = tx.Where("spd.size_id = ? AND spd.remaining > ? AND spd.status = ?", sizeId, 0, "active").Where("pd.id = ?", productItemDetailId).Where("p.id = ?", producId).Where("s.id = ?", sizeId)

	// tx = tx.Preload("StockProducts", func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("remaining > ? AND status = ?", 0, "active")
	// })

	tx = tx.Group("pd.product_id,pd.id,p.name, s.size, pd.color_name, s.price,s.id")

	err := tx.Find(&products).Error
	return products, err
}

func GetProductByIDSizdID(db *gorm.DB, productItemDetailId, sizeId, producId string) (custommodel.ProductOrderDetails, error) {
	var products custommodel.ProductOrderDetails

	tx := db.Table(models.TableNameProductDetail + " pd")

	tx = tx.Select("pd.product_id, pd.id, p.name, s.id  AS size_id , s.size ,pd.color_name , pd.price")

	tx = tx.Where("pd.id = ?", productItemDetailId).Where("p.id = ?", producId).Where("s.id = ?", sizeId)
	// .Where("spd.status = ?", "active")

	joinProduct := fmt.Sprintf("LEFT JOIN %s p ON pd.product_id = p.id", models.TableNameProduct)
	joinQuery := fmt.Sprintf("LEFT JOIN %s spd ON spd.product_detail_id = pd.id", models.TableNameStockProductDetail)
	joinSize := fmt.Sprintf("LEFT JOIN %s s ON s.product_detail_id = pd.id", models.TableNameSize)

	tx = tx.Joins(joinQuery).Joins(joinProduct).Joins(joinSize)

	// tx = tx.Where("spd.size_id = ? AND spd.remaining > ? AND spd.status = ?", sizeId, 0, "active").Where("pd.id = ?", productItemDetailId).Where("p.id = ?", producId).Where("s.id = ?", sizeId)

	// tx = tx.Preload("StockProducts", func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("remaining > ? AND status = ?", 0, "active")
	// })

	tx = tx.Group("pd.product_id,pd.id,p.name, s.size, pd.color_name, pd.price,s.id")

	err := tx.Find(&products).Error
	return products, err
}
