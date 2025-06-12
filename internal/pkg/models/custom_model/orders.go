package custommodel

import (
	"go-api/internal/pkg/models"
	"time"
)

// OrderRequest represents the order payload with validation tags
type OrderRequest struct {
	CustomAddress string `json:"custom_address" validate:"required"`
	FullName      string `json:"full_name" validate:"required"`
	Tel           int32  `json:"tel" validate:"required,min=8"`
	PlatForm      string `json:"plat_form" validate:"required,oneof=facebook whatapp tiktok"`
	Gender        int    `json:"gender" validate:"required,oneof=0 1 2"` // 0: other, 1: male, 2: female
	Cod           bool   `json:"cod" `
	FreeShipping  bool   `json:"free_shipping" `
	// DiscountStr   string      `json:"discount_" validate:"omitempty,numeric"`
	// ProvinceID int    `json:"province_id" validate:"required,min=1"`
	ShippingID string `json:"shipping_id" validate:"required"`
	DistrictID int32  `json:"district_id" validate:"required,min=1"`
	// TotalPrice float64     `json:"total_price" validate:"required,min=0"`
	Discount   float64     `json:"discount" validate:"min=0"`
	FacebookID string      `json:"facebook_id" validate:"required"`
	Items      []OrderItem `json:"items" validate:"required,dive,required"`
}

// OrderItem represents each product in the order
type OrderItem struct {
	ProductID       string `json:"product_id" validate:"required"`
	ProductDetailID string `json:"product_detail_id" validate:"required"`
	SizeID          string `json:"size_id" validate:"required"`
	Quantity        int32  `json:"quantity" validate:"required,min=1"`
}

// Original input item
type Item struct {
	ProductID       string `json:"product_id"`
	ProductDetailID string `json:"product_detail_id"`
	SizeID          string `json:"size_id"`
	Quantity        int    `json:"quantity"`
}

// Grouped item without product_id
type GroupedItem struct {
	ProductDetailID string `json:"product_detail_id"`
	SizeID          string `json:"size_id"`
	Quantity        int32  `json:"quantity"`
	ProductDetails  ProductOrderDetails
}

// Final output with total quantities
type GroupedByProduct struct {
	ProductID           string  `json:"product_id"`
	PromotionID         string  `json:"promotion_id"`
	DiscountOnlyProduct float32 `json:"discount_only_product"`
	Promotion           []models.Promotion
	TotalQuantities     int32         `json:"total_quantities"`
	TotalPrice          float64       `json:"total_prices"`
	Items               []GroupedItem `json:"items"`
}

type OrderReponse struct {
	ID                    string          `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Status                string          `gorm:"column:status;default:pending" json:"status"`
	CustomerID            string          `gorm:"column:customer_id" json:"customer_id"`
	IsCancel              bool            `gorm:"column:is_cancel;not null" json:"is_cancel"`
	PageName              string          `gorm:"column:page_name" json:"page_name"`
	Tel                   int32           `gorm:"column:tel" json:"tel"`
	CustomAddress         string          `gorm:"column:custom_address" json:"custom_address"`
	TotalProductsDiscount float64         `gorm:"column:total_prodouct_discount" json:"total_prodouct_discount"`
	UserID                string          `gorm:"column:user_id" json:"user_id"`
	District              string          `gorm:"column:dr_name" json:"dr_name"`
	Province              string          `gorm:"column:pr_name" json:"pr_name"`
	TotalPrice            float64         `gorm:"column:total_price" json:"total_price"`
	DistrictID            int32           `gorm:"column:district_id" json:"district_id"`
	OrderedAt             time.Time       `gorm:"column:ordered_at;default:now()" json:"ordered_at"`
	UpdatedAt             time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	OrderNo               string          `gorm:"column:order_no;not null" json:"order_no"`
	FreeShipping          bool            `gorm:"column:free_shipping;not null" json:"free_shipping"`
	OrderName             string          `gorm:"column:order_name;not null" json:"order_name"`
	ShippingID            string          `gorm:"column:shipping_id;not null" json:"shipping_id"`
	CreatedAt             time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Discount              float64         `gorm:"column:discount;not null" json:"discount"`
	Platform              string          `gorm:"column:platform;default:facebook" json:"platform"`
	Cod                   bool            `gorm:"column:cod;default:true" json:"cod"`
	OrderDetails          []OrderDetail   `gorm:"foreignKey:OrderID;references:ID" json:"order_details"`
	Shipping              models.Shipping `gorm:"foreignKey:ShippingID;references:ID" json:"shipping"`
	OrderDiscounts        []OrderDiscount `gorm:"foreignKey:OrderID;references:ID" json:"order_discounts"`
}

// // TableName Order's table name
// func (*OrderReponse) TableName() string {
// 	return models.TableNameOrder
// }

// OrderDetail mapped from table <order_details>
type OrderDetail struct {
	ID              string        `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID         string        `gorm:"column:order_id;not null" json:"order_id"`
	ProductDetailID string        `gorm:"column:product_detail_id;not null" json:"product_detail_id"`
	Quantity        int32         `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice       float64       `gorm:"column:unit_price;not null" json:"unit_price"`
	TotalPrice      float64       `gorm:"column:total_price;not null" json:"total_price"`
	CreatedAt       time.Time     `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	SizeID          string        `gorm:"column:size_id;not null" json:"size_id"`
	ProductDetail   ProductDetail `gorm:"foreignKey:ProductDetailID;references:ID" json:"product_detail"`
	Size            SizeOrder     `gorm:"foreignKey:SizeID;references:ID" json:"size"`
}

// TableName OrderDetail's table name
func (*OrderDetail) TableName() string {
	return models.TableNameOrderDetail
}

// OrderDiscount mapped from table <order_discounts>
type OrderDiscount struct {
	ID            string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID       string    `gorm:"column:order_id;not null" json:"order_id"`
	ProductID     string    `gorm:"column:product_id;not null" json:"product_id"`
	TotalDiscount float64   `gorm:"column:total_discount;not null" json:"total_discount"`
	DiscountID    string    `gorm:"column:discount_id;not null" json:"discount_id"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName OrderDiscount's table name
func (*OrderDiscount) TableName() string {
	return models.TableNameOrderDiscount
}

type ProductDetail struct {
	ID        string `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID string `gorm:"column:product_id;not null" json:"product_id"`
	Color     string `gorm:"column:color" json:"color"`
	FitType   string `gorm:"column:fit_type" json:"fit_type"`
	Material  string `gorm:"column:material" json:"material"`
	// Status    string    `gorm:"column:status;default:active" json:"status"`
	// CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	// UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ImageURL string `gorm:"column:image_url;default:N/A" json:"image_url"`
	// Price     int32   `gorm:"column:price;not null" json:"price"`
	ColorName string  `gorm:"column:color_name" json:"color_name"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// TableName ProductDetail's table name
func (*ProductDetail) TableName() string {
	return models.TableNameProductDetail
}

type SizeOrder struct {
	ID              string  `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Size            string  `gorm:"column:size" json:"size"`
	ProductDetailID string  `gorm:"column:product_detail_id;not null" json:"-"`
	Price           float64 `gorm:"column:price;not null" json:"price"`
}

// TableName Size's table name
func (*SizeOrder) TableName() string {
	return models.TableNameSize
}

// Product mapped from table <products>
type Product struct {
	ID    string `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"-"`
	Name  string `gorm:"column:name;not null" json:"name"`
	Brand string `gorm:"column:brand" json:"brand"`
	// CategoryID string    `gorm:"column:category_id;not null" json:"category_id"`
	// Price      float64   `gorm:"column:price;not null" json:"price"`
	// Status     string    `gorm:"column:status;default:active" json:"status"`
	// CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	// UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Product's table name
func (*Product) TableName() string {
	return models.TableNameProduct
}
