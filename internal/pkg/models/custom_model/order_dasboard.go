package custommodel

type FlatProductSales struct {
	ProductName string  `json:"product_name" gorm:"column:product_name"`
	Hour        string  `json:"hour" gorm:"column:hour"`
	Quantity    int     `json:"quantity" gorm:"column:total_quantity"`
	TotalPrice  float64 `json:"total_price" gorm:"column:total_price"`
}

type HourlyStat struct {
	Hour       string  `json:"hour"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type ProductSalesGrouped struct {
	ProductName string       `json:"product_name"`
	Times       []HourlyStat `json:"times"`
}
