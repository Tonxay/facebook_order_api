package dbservice

import (
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
)

func CreateMesseng(messeng *models.Chat) error {

	err := gormpkg.GetDB().Table(models.TableNameChat).Create(&messeng).Error

	return err
}

func GetMessengerPerUser(userId string, pageId string) ([]models.Chat, error) {
	var messages []models.Chat
	err := gormpkg.GetDB().
		Table(models.TableNameChat).
		Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			userId, pageId, pageId, userId).
		Find(&messages).Error
	return messages, err
}
