package custommodel

import "go-api/internal/pkg/models"

// OrderRequest represents the order payload with validation tags
type OrderRequest struct {
	CustomAddress string `json:"custom_address" validate:"required"`
	FullName      string `json:"full_name" validate:"required"`
	Tel           string `json:"tel" validate:"required,e164|numeric,min=8,max=20"`
	PhatForm      string `json:"phat_form" validate:"required,oneof=facebook whatapp tiktok"`
	Gender        int    `json:"gender" validate:"required,oneof=0 1 2"` // 0: other, 1: male, 2: female
	PayType       bool   `json:"pay_type"`
	ShippingCost  bool   `json:"shipping_cost" `
	// DiscountStr   string      `json:"discount_" validate:"omitempty,numeric"`
	ProvinceID int    `json:"province_id" validate:"required,min=1"`
	Shipping   string `json:"shipping" validate:"required"`
	DistrictID int    `json:"district_id" validate:"required,min=1"`
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
	Quantity        int    `json:"quantity" validate:"required,min=1"`
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
	Quantity        int    `json:"quantity"`
	ProductDetails  ProductOrderDetails
}

// Final output with total quantities
type GroupedByProduct struct {
	ProductID           string `json:"product_id"`
	DiscountOnlyProduct int    `json:"discount_only_product"`
	Promotion           []models.Promotion
	TotalQuantities     int           `json:"total_quantities"`
	TotalPrice          int32         `json:"total_prices"`
	Items               []GroupedItem `json:"items"`
}
