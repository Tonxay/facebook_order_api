package custommodel

import (
	"go-api/internal/pkg/models"
)

// Product mapped from table <products>
type Products struct {
	ID             string           `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string           `gorm:"column:name;not null" json:"name"`
	ProductDetails []ProductDetails `gorm:"foreignKey:ProductID;references:ID" json:"product_details"`
	Promotions     []Promotion      `gorm:"foreignKey:ProductID;references:ID" json:"promotions"`
}

// TableName Product's table name
func (*Products) TableName() string {
	return models.TableNameProduct
}

type ProductDetails struct {
	ID        string `gorm:"column:id" json:"id"`
	ProductID string `gorm:"column:product_id;not null" json:"product_id"`
	Color     string `gorm:"column:color" json:"color"`
	FitType   string `gorm:"column:fit_type" json:"fit_type"`
	Material  string `gorm:"column:material" json:"material"`
	Status    string `gorm:"column:status;default:active" json:"status"`
	ImageURL  string `gorm:"column:image_url;default:N/A" json:"image_url"`

	// Price     int32  `gorm:"column:price;not null" json:"price"`
	ColorName string `gorm:"column:color_name" json:"color_name"`
	Sizes     []Size `gorm:"foreignKey:product_detail_id;references:ID" json:"sizes"`
}

// TableName ProductDetail's table name
func (*ProductDetails) TableName() string {
	return models.TableNameProductDetail
}

type Size struct {
	ID              string  `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Size            string  `gorm:"column:size" json:"size"`
	TotalRemaining  int32   `gorm:"column:total_remaining" json:"total_remaining"`
	Price           float64 `gorm:"column:price;not null" json:"price"`
	ProductDetailID string  `gorm:"column:product_detail_id;not null" json:"product_detail_id"`
}

// TableName Size's table name
func (*Size) TableName() string {
	return models.TableNameSize
}

type Promotion struct {
	ID        string  `gorm:"column:id;not null;default:gen_random_uuid()" json:"id"`
	ProductID string  `gorm:"column:product_id;not null" json:"product_id"`
	Discount  float32 `gorm:"column:discount;not null" json:"discount"`
	Status    string  `gorm:"column:status;not null;default:active" json:"status"`
	Quentity  int32   `gorm:"column:quentity;not null" json:"quentity"`
}

// TableName Promotion's table name
func (*Promotion) TableName() string {
	return models.TableNamePromotion
}

type ProductOrderDetails struct {
	ID        string `gorm:"column:id" json:"id"`
	ProductID string `gorm:"column:product_id;not null" json:"product_id"`
	Name      string `gorm:"column:name;not null" json:"name"`
	// Color     string `gorm:"column:color" json:"color"`
	// FitType   string `gorm:"column:fit_type" json:"fit_type"`
	// Material  string `gorm:"column:material" json:"material"`
	// Status    string `gorm:"column:status;default:active" json:"status"`
	Price     float64 `gorm:"column:price;not null" json:"price"`
	ColorName string  `gorm:"column:color_name" json:"color_name"`
	Remaining int32   `gorm:"column:remaining;not null" json:"remaining"`
	SizeID    string  `gorm:"column:size_id" json:"size_id"`
	Size      string  `gorm:"column:size" json:"size"`
}
