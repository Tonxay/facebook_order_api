package dbservice

import (
	"context"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/query"

	"gorm.io/gorm"
)

func CreateStockProductDetail(db *gorm.DB, stockProductDetail *models.StockProductDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.StockProductDetail
	err := daq.WithContext(ctx).Create(stockProductDetail)
	return err
}

func CreateStockProductDetailForOrder(db *gorm.DB, stockProductDetail []*models.StockProductDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.StockProductDetail
	err := daq.WithContext(ctx).CreateInBatches(stockProductDetail, 100)
	return err
}
func CreateRoderStockDetails(db *gorm.DB, stockDetails []*models.OrderStockDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.OrderStockDetail
	err := daq.WithContext(ctx).CreateInBatches(stockDetails, 100)
	return err
}

func UpdateStockProductDetail(db *gorm.DB, id string, remaining int32, status string) error {
	err := db.Table(models.TableNameStockProductDetail).
		Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"remaining": remaining,
			"status":    status,
		}).Error

	return err
}

// func SubtractRemaining(db *gorm.DB, productDetailID string, quantity int32) error {
// 	result := db.Model(&ProductDetail{}).
// 		Where("id = ? AND remaining >= ?", productDetailID, quantity).
// 		Update("remaining", gorm.Expr("remaining - ?", quantity))

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return fmt.Errorf("not enough stock or product_detail_id not found")
// 	}

// 	return nil
// }
