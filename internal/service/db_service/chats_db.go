package dbservice

import (
	gormpkg "github.com/yourusername/go-api/internal/pkg"
	"github.com/yourusername/go-api/internal/pkg/models"
)

func CreateMesseng(messeng models.Chat) error {

	err := gormpkg.GetDB().Table(models.TableNameChat).Create(&messeng).Error
	return err
}
