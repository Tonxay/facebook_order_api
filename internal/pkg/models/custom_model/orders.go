package custommodel

type OrderRequest struct {
	CustomerID    string         `gorm:"column:customer_id" json:"customer_id"`
	Tel           int32          `gorm:"column:tel" json:"tel"`
	CustomAddress string         `gorm:"column:custom_address" json:"custom_address"`
	UserID        string         `gorm:"column:user_id" json:"-"`
	TotalPrice    int32          `gorm:"column:total_price" json:"-"`
	DistrictID    int32          `gorm:"column:district_id" json:"district_id"`
	OrderDetails  []OrderDetails `json:"items"`
}

type OrderDetails struct {
	ProductDetailID string `gorm:"column:product_detail_id;not null" json:"product_detail_id"`
	Quantity        int32  `gorm:"column:quantity;not null" json:"quantity"`
	SizeID          string `gorm:"column:size_id;not null" json:"size_id"`
}
