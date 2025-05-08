package dbservice

import (
	"github.com/yourusername/go-api/internal/pkg/models"
	"github.com/yourusername/go-api/internal/pkg/query"
)

func CreateMesseng(messeng models.Chat) error {
	err := query.Chat.Create(&messeng)
	return err
}
