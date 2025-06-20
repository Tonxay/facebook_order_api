package request

type ProductManyRequest struct {
	ProductName    string          `json:"product_name" validate:"required"`
	Brand          string          `json:"brand" validate:"required"`
	CategoryID     string          `json:"category_id" validate:"required,uuid4"`
	ProductDetails []ProductDetail `json:"product_details" validate:"required,dive"`
}

type ProductDetail struct {
	ColorName string        `json:"color_name" validate:"required"`
	Color     string        `json:"color" validate:"required,hexcolor|eq=N/A"`
	Material  string        `json:"material" validate:"required"`
	FitType   string        `json:"fit_type" validate:"required"`
	ImageURL  string        `json:"image_url" validate:"required,url|eq=N/A"`
	Sizes     []ProductSize `json:"sizes" validate:"required,dive"`
}

type ProductSize struct {
	Size  string  `json:"size" validate:"required"`
	Price float64 `json:"price" validate:"required,gt=0"`
}
