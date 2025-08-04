package custommodel

type FilterDasboard struct {
	ProductID string `query:"product_id" validate:"required,uuid4" json:"product_id"`
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02" json:"start_date"`
	EndDate   string `query:"end_date" validate:"omitempty,datetime=2006-01-02" json:"end_date"`
}
