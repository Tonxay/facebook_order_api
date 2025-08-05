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

type OrderSummary struct {
	TotalOrders        int     `json:"total_orders"`
	Female             int     `json:"female"`
	Male               int     `json:"male"`
	SalesRepsInvolved  int     `json:"sales_reps_involved"`
	TotalUnitsSold     float64 `json:"total_units_sold"`
	GrossRevenue       float64 `json:"gross_revenue"`
	TotalDiscounts     float64 `json:"total_discounts"`
	NetRevenue         float64 `json:"net_revenue"`
	AvgOrderValue      float64 `json:"avg_order_value"`
	FreeShippingOrders int     `json:"free_shipping_orders"`
	CodOrders          int     `json:"cod_orders"`
	Paymented          int     `json:"paymented"`
}
