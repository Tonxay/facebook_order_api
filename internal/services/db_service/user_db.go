package dbservice

import (
	"context"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/query"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User, ctx context.Context) error {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	query.SetDefault(db)
	daq := query.Q.User
	err := daq.WithContext(ctx).Create(user)
	return err
}

func GetUserForUserName(db *gorm.DB, userName string) (custommodel.User, error) {
	var user custommodel.User
	err := db.Table(models.TableNameUser).Where("user_name = ?", userName).First(&user).Error
	return user, err
}

func GetUserNamePassword(db *gorm.DB, userName string) (models.User, error) {
	var user models.User
	err := db.Table(models.TableNameUser).Where("user_name = ?", userName).First(&user).Error
	return user, err
}

func GetUserBayID(db *gorm.DB, id string) (custommodel.User, error) {
	var user custommodel.User
	err := db.Table(models.TableNameUser).Where("id = ?", id).First(&user).Error
	return user, err
}
