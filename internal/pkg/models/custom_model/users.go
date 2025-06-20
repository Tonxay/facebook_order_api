package custommodel

import "go-api/internal/pkg/models"

// User mapped from table <users>
type User struct {
	ID       string `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Password string `gorm:"column:password;not null" json:"-"`
	UserName string `gorm:"column:user_name;not null" json:"user_name"`
	Status   string `gorm:"column:status;not null;default:active" json:"status"`
	PageID   string `gorm:"column:page_id;not null" json:"page_id"`
}

// TableName User's table name
func (*User) TableName() string {
	return models.TableNameUser
}

// Customer mapped from table <customers>
type Customer struct {
	FacebookID  string      `gorm:"column:facebook_id;primaryKey" json:"facebook_id"`
	LastName    string      `gorm:"column:last_name" json:"last_name"`
	PhoneNumber int64       `gorm:"column:phone_number" json:"phone_number"`
	FirstName   string      `gorm:"column:first_name" json:"first_name"`
	PageID      string      `gorm:"column:page_id;not null" json:"page_id"`
	NamePage    string      `gorm:"column:name_page;not null" json:"name_page"`
	Gender      int32       `gorm:"column:gender;not null" json:"gender"`
	Page        models.Page `gorm:"foreignKey:PageID;references:PageID" json:"page"`
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return models.TableNamePage
}
