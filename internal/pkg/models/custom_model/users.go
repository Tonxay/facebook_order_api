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
