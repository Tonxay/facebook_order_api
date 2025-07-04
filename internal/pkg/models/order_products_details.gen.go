// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameOrderProductsDetail = "order_products_details"

// OrderProductsDetail mapped from table <order_products_details>
type OrderProductsDetail struct {
	ID              string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	ProductDetailID string    `gorm:"column:product_detail_id;not null" json:"product_detail_id"`
	UnitPrice       float64   `gorm:"column:unit_price;not null" json:"unit_price"`
	TotalPrice      float64   `gorm:"column:total_price;not null" json:"total_price"`
	Quantity        int32     `gorm:"column:quantity;not null" json:"quantity"`
	SizeID          string    `gorm:"column:size_id;not null" json:"size_id"`
	OrderProductID  string    `gorm:"column:order_product_id" json:"order_product_id"`
}

// TableName OrderProductsDetail's table name
func (*OrderProductsDetail) TableName() string {
	return TableNameOrderProductsDetail
}
