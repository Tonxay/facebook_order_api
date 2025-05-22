package custommodel

import "go-api/internal/pkg/models"

type ProductDetailCounter struct {
	ID            string                      `gorm:"column:id;not null" json:"id"`
	Name          string                      `gorm:"column:name;not null" json:"name"`
	Size          string                      `gorm:"column:size" json:"size"`
	Color         string                      `gorm:"column:color" json:"color"`
	Quantity      int32                       `gorm:"column:quantity;not null" json:"quantity"`
	Price         int32                       `gorm:"column:price;not null" json:"price"`
	StockProducts []models.StockProductDetail `gorm:"foreignKey:ProductDetailID;references:ID" json:"stock_product"`
}

// TableName ProductDetail's table name
func (*ProductDetailCounter) TableName() string {
	return models.TableNameProductDetail
}

type ProductDetails struct {
	ID       string `gorm:"column:id;not null" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Size     string `gorm:"column:size" json:"size"`
	Color    string `gorm:"column:color" json:"color"`
	Quantity int32  `gorm:"column:quantity;not null" json:"quantity"`
	Price    int32  `gorm:"column:price;not null" json:"price"`
	// StockProducts []models.StockProductDetail `gorm:"foreignKey:ProductDetailID;references:ID" json:"stock_product"`
}

// TableName ProductDetail's table name
func (*ProductDetails) TableName() string {
	return models.TableNameProductDetail
}
