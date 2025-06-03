package request

type User struct {
	Password string `gorm:"column:password;not null" validate:"required,min=6" json:"password"`
	UserName string `gorm:"column:user_name;not null" validate:"email,required" json:"user_name"`
}

type CustomerQuery struct {
	StartDate string `query:"start_date" validate:"required"`
	EndDate   string `query:"end_date" validate:"required"`
	Name      string `query:"name"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
