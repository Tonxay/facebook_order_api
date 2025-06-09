package custommodel

import "go-api/internal/pkg/models"

// OrderRequest represents the order payload with validation tags
type OrderRequest struct {
	// CustomAddress string `json:"custom_address" validate:"required"`
	FullName     string `json:"full_name" validate:"required"`
	Tel          int32  `json:"tel" validate:"required,min=8"`
	PlatForm     string `json:"plat_form" validate:"required,oneof=facebook whatapp tiktok"`
	Gender       int    `json:"gender" validate:"required,oneof=0 1 2"` // 0: other, 1: male, 2: female
	Cod          bool   `json:"cod" `
	FreeShipping bool   `json:"free_shipping" `
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
