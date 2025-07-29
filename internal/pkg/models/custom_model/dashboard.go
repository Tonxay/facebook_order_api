package custommodel

type ProductForProvincesReport struct {
	ProvinceName      string  `json:"province_name"`
	ProvinceNameEn    string  `json:"province_name_en"`
	ProductName       string  `json:"product_name"`
	OrderCount        int64   `json:"order_count"`
	TotalQuantity     float64 `json:"total_quantity"`
	PromotionDiscount float64 `json:"promotion_discount"`
	ProductDiscount   float64 `json:"product_discount"`
	TotalProductPrice float64 `json:"total_product_price"`
}
