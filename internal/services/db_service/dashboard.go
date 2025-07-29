package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"log"

	"gorm.io/gorm"
)

func GetProductforProvince(db *gorm.DB, productID string) ([]custommodel.ProductForProvincesReport, error) {
	var reports []custommodel.ProductForProvincesReport
	tx := db.Table(models.TableNameProvince + " pr").Select(
		`
						pr.pr_name AS province_name,
						pr.pr_name_en AS province_name_en,
						p.name AS product_name,
						COUNT(*) AS order_count,
				        SUM(op.total_product_price) AS total_product_price,
						SUM(opd.discount) AS promotion_discount,
						SUM(op.discount) AS  product_discount,
						SUM(op.total_amounts) AS total_quantity
						`,
	)
	tx =
		tx.
			Joins("JOIN "+models.TableNameDistrict+" d ON d.province_id = pr.id").
			Joins("JOIN "+models.TableNameOrder+" o ON o.district_id = d.id").
			Joins("JOIN "+models.TableNameOrderProduct+" op ON o.id = op.order_id").
			Joins("LEFT JOIN "+models.TableNameOrderProductDiscount+" opd ON opd.order_product_id = op.id").
			Joins("JOIN products p ON op.product_id = p.id").
			Where("o.is_cancel IS NOT TRUE AND p.id = ?", productID)

	err := tx.Group(`pr.pr_name,pr.pr_name_en,p.name`).Order("total_product_price DESC").
		Scan(&reports).Error

	if err != nil {
		log.Println("query error:", err)
	}
	return reports, err
}
