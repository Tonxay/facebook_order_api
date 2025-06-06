// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameOrder = "orders"

// Order mapped from table <orders>
type Order struct {
	ID            string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Status        string    `gorm:"column:status;default:pending" json:"status"`
	CustomerID    string    `gorm:"column:customer_id" json:"customer_id"`
	Tel           int32     `gorm:"column:tel" json:"tel"`
	CustomAddress string    `gorm:"column:custom_address" json:"custom_address"`
	UserID        string    `gorm:"column:user_id" json:"user_id"`
	TotalPrice    int32     `gorm:"column:total_price" json:"total_price"`
	DistrictID    int32     `gorm:"column:district_id" json:"district_id"`
	OrderedAt     time.Time `gorm:"column:ordered_at;default:CURRENT_TIMESTAMP" json:"ordered_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	PackagePrice  int32     `gorm:"column:package_price" json:"package_price"`
	OrderNo       string    `gorm:"column:order_no;not null" json:"order_no"`
	PayType       bool      `gorm:"column:pay_type;not null" json:"pay_type"`
	OrderName     int32     `gorm:"column:order_name" json:"order_name"`
	ShippingID    string    `gorm:"column:shipping_id;not null" json:"shipping_id"`
}

// TableName Order's table name
func (*Order) TableName() string {
	return TableNameOrder
}
