package dbservice

import (
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
)

func CreateMesseng(messeng *models.Chat) error {

	err := gormpkg.GetDB().Table(models.TableNameChat).Create(&messeng).Error
	return err
}
