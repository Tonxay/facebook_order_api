package request

type User struct {
	Password string `gorm:"column:password;not null" validate:"required,min=6" json:"password"`
	UserName string `gorm:"column:user_name;not null" validate:"email,required" json:"user_name"`
}
