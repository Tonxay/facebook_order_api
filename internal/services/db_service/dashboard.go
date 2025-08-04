package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"log"

	"gorm.io/gorm"
)

func GetProductforProvince(db *gorm.DB, filter custommodel.FilterDasboard) ([]custommodel.ProductForProvincesReport, error) {
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
			Joins("JOIN "+models.TableNameProduct+" p ON op.product_id = p.id").
			Where("o.is_cancel IS NOT TRUE AND p.id = ?", filter.ProductID).
			Where("o.ordered_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)

	err := tx.Group(`pr.pr_name,pr.pr_name_en,p.name`).Order("total_product_price DESC").
		Scan(&reports).Error

	if err != nil {
		log.Println("query error:", err)
	}
	return reports, err
}

func GetOrderSummary(db *gorm.DB, filter custommodel.FilterDasboard) (custommodel.OrderSummary, error) {
	var summary custommodel.OrderSummary

	tx := db.
		Table("orders AS o").
		Select(`
			COUNT(DISTINCT o.id) AS total_orders,
			COUNT(DISTINCT o.customer_id) AS unique_customers,
			COUNT(DISTINCT o.user_id) AS sales_reps_involved,
			SUM(op.total_amounts) AS total_units_sold,
			SUM(op.total_product_price) AS gross_revenue,
			SUM(opd.discount + op.discount) AS total_discounts,
			(SUM(op.total_product_price) - SUM(opd.discount + op.discount)) AS net_revenue,
			CAST(AVG(op.total_product_price) AS DECIMAL(10,2)) AS avg_order_value,
			COUNT(CASE WHEN o.free_shipping THEN 1 END) AS free_shipping_orders,
			COUNT(CASE WHEN o.cod THEN 1 END) AS cod_orders
		`).
		Joins("LEFT JOIN order_products op ON o.id = op.order_id").
		Joins("LEFT JOIN order_product_discounts opd ON op.id = opd.order_product_id").
		Where("o.is_cancel = ?", false)

	if filter.StartDate == "" {
		tx = tx.Where("o.ordered_at = CURRENT_DATE")
	} else {
		tx = tx.Where("o.ordered_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	err := tx.Scan(&summary).Error

	if err != nil {
		return summary, err
	}

	return summary, nil
}
