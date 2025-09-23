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
			COUNT(o.id) AS total_orders,
			COUNT(CASE WHEN c.gender = 2 THEN  1 END ) as female,
            COUNT(CASE WHEN c.gender = 1 THEN  1 END ) as male,
			SUM(op.total_amounts) AS total_units_sold,
			SUM(op.total_product_price) AS gross_revenue,
			SUM(opd.discount + op.discount) AS total_discounts,
			(SUM(op.total_product_price) - SUM(opd.discount + op.discount)) AS net_revenue,
			CAST(AVG(op.total_product_price) AS DECIMAL(10,2)) AS avg_order_value,
			COUNT(CASE WHEN o.free_shipping THEN 1 END) AS free_shipping_orders,
			COUNT(CASE WHEN o.cod THEN 1 END) AS cod_orders,
			COUNT(CASE WHEN o.cod = false THEN 1 END) as paymented
		`).
		Joins("LEFT JOIN order_products op ON o.id = op.order_id").
		Joins("LEFT JOIN customers c ON c.facebook_id = o.customer_id").
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
func GetProductSalesByDay(db *gorm.DB, filter custommodel.FilterDasboard) ([]custommodel.ProductSalesByDay, error) {
	var results []custommodel.ProductSalesByDay

	// Build the base query using GORM
	query := db.Table("orders o").
		Select(`
			CASE EXTRACT(dow FROM o.ordered_at) 
				WHEN 0 THEN 'Sunday' 
				WHEN 1 THEN 'Monday' 
				WHEN 2 THEN 'Tuesday' 
				WHEN 3 THEN 'Wednesday' 
				WHEN 4 THEN 'Thursday' 
				WHEN 5 THEN 'Friday' 
				WHEN 6 THEN 'Saturday' 
			END as day_of_week,
			EXTRACT(dow FROM o.ordered_at) as dow_number,
			p.name as product_name,
			p.brand,
			COUNT(DISTINCT o.id) as total_orders,
			SUM(op.total_amounts) as total_units_sold,
			SUM(op.total_product_price) as total_revenue,
			AVG(op.total_amounts) as avg_units_per_order
		`).
		Joins("JOIN order_products op ON o.id = op.order_id").
		Joins("JOIN products p ON op.product_id = p.id").
		Joins("LEFT JOIN order_product_discounts opdc ON op.id = opdc.order_product_id").
		Where("o.ordered_at IS NOT NULL").
		Where("o.is_cancel = ?", false)

	// Apply filters
	if filter.ProductID != "" {
		query = query.Where("p.id = ?", filter.ProductID)
	}

	if filter.StartDate != "" {
		query = query.Where("DATE(o.ordered_at) >= ?", filter.StartDate)
	}

	if filter.EndDate != "" {
		query = query.Where("DATE(o.ordered_at) <= ?", filter.EndDate)
	}

	// Group by and order
	query = query.
		Group("EXTRACT(dow FROM o.ordered_at), p.id, p.name, p.brand").
		Order("dow_number, total_units_sold DESC")

	// Execute the query
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func GetProductOrderCount(db *gorm.DB, filter custommodel.FilterDasboard) ([]custommodel.ProductOrderCount, error) {
	var results []custommodel.ProductOrderCount
	err := db.Table("orders AS o").
		Select(`
        p.name AS product_name,
		p.id AS product_id,
        COUNT(op.id) AS order_count,
        o.ordered_at AS day
    `).
		Joins("JOIN order_products op ON o.id = op.order_id").
		Joins("JOIN products p ON op.product_id = p.id").
		Where("o.ordered_at IS NOT NULL").
		Where("o.ordered_at >= ?", filter.StartDate).
		Where("o.ordered_at <= ?", filter.EndDate).
		Where("o.is_cancel IS NOT TRUE").
		Group("p.name, o.ordered_at,p.id").
		Order("o.ordered_at, order_count DESC").
		Scan(&results).Error

	if err != nil {
		log.Println("Error:", err)
	}
	return results, err
}

func GetProductSales(db *gorm.DB, filter custommodel.FilterDasboard) ([]custommodel.ProductSales, error) {
	var results []custommodel.ProductSales

	query := db.Model(&models.Order{}).
		Select(`
        products.name AS name,
        COUNT(DISTINCT orders.id) AS total_orders,
        SUM(op.total_product_price - op.discount - COALESCE(order_product_discounts.discount, 0)) AS net,
     	SUM(op.total_amounts) as total_units_sold,
        SUM(SUM(op.total_product_price - op.discount - COALESCE(order_product_discounts.discount, 0)))
            OVER () AS total_net
    `).
		Joins("LEFT JOIN order_products op ON orders.id = op.order_id").
		Joins("LEFT JOIN products ON op.product_id = products.id").
		Joins("LEFT JOIN order_product_discounts ON op.id = order_product_discounts.order_product_id").
		Where("orders.is_cancel = ?", false).
		Where("orders.ordered_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate).
		Order("total_units_sold ASC").
		Group("products.name")

	err := query.Scan(&results).Error
	return results, err
}
