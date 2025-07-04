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
	CustomerID    string    `gorm:"column:customer_id;not null" json:"customer_id"`
	Tel           int64     `gorm:"column:tel" json:"tel"`
	CustomAddress string    `gorm:"column:custom_address" json:"custom_address"`
	UserID        string    `gorm:"column:user_id" json:"user_id"`
	TotalPrice    float64   `gorm:"column:total_price" json:"total_price"`
	DistrictID    int32     `gorm:"column:district_id" json:"district_id"`
	OrderedAt     time.Time `gorm:"column:ordered_at;default:now()" json:"ordered_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	OrderNo       string    `gorm:"column:order_no;not null" json:"order_no"`
	FreeShipping  bool      `gorm:"column:free_shipping;not null" json:"free_shipping"`
	OrderName     string    `gorm:"column:order_name;not null" json:"order_name"`
	ShippingID    string    `gorm:"column:shipping_id;not null" json:"shipping_id"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Discount      float64   `gorm:"column:discount;not null" json:"discount"`
	Platform      string    `gorm:"column:platform;default:facebook" json:"platform"`
	Cod           bool      `gorm:"column:cod;not null" json:"cod"`
	IsCancel      bool      `gorm:"column:is_cancel;not null" json:"is_cancel"`
	UserUpdated   string    `gorm:"column:user_updated;not null" json:"user_updated"`
}

// TableName Order's table name
func (*Order) TableName() string {
	return TableNameOrder
}
