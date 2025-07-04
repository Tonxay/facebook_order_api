package request

import "go-api/internal/pkg/models"

type StockProductDetail struct {
	ID              string `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductDetailID string `gorm:"column:product_detail_id;not null" json:"product_detail_id"`
	Quantity        int32  `gorm:"column:quantity;not null" json:"quantity"`
	Status          string `gorm:"column:status;default:active" json:"status"`
	SizeID          string `gorm:"column:size_id" json:"size_id"`
	UserID          string `gorm:"column:user_id" json:"user_id"`
	Remaining       int32  `gorm:"column:remaining" json:"-"`
}

// TableName StockProductDetail's table name
func (*StockProductDetail) TableName() string {
	return models.TableNameStockProductDetail
}

type StockIncreaseRequest struct {
	StockIncrease []StockIncrease `json:"items" validate:"required,dive"`
}

type StockIncrease struct {
	ProductDetailID string `json:"product_detail_id"`
	SizeID          string `json:"size_id"`
	Quantity        int32  `json:"quantity"`
	Remaining       int32  `json:"remaining"`
}
