package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	gormpkg "github.com/yourusername/go-api/internal/pkg"
	"github.com/yourusername/go-api/internal/pkg/models"
	"github.com/yourusername/go-api/internal/pkg/models/customs"
	dbservice "github.com/yourusername/go-api/internal/service/db_service"
)

type WebhookDeliveryEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`

			Timestamp float64 `json:"timestamp"`

			Message *struct {
				Mid         string `json:"mid"`
				Text        string `json:"text,omitempty"`
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						URL string `json:"url"`
					} `json:"payload"`
				} `json:"attachments,omitempty"`
			} `json:"message,omitempty"`

			Delivery *struct {
				Mids      []string `json:"mids"`
				Watermark float64  `json:"watermark"`
			} `json:"delivery,omitempty"`

			Postback *struct {
				Payload string `json:"payload"`
			} `json:"postback,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}

type WebhookEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		Messaging []struct {
			Sender  struct{ ID string } `json:"sender"`
			Message struct {
				Text string `json:"text"`
			} `json:"message,omitempty"`
			Postback struct {
				Payload string `json:"payload"`
			} `json:"postback,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}

func VerifyWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == os.Getenv("VERIFY_TOKEN") {
		return c.SendString(challenge)
	}
	return c.SendStatus(fiber.StatusForbidden)
}

// Raw Event: map[entry:[map[changes:[map[field:inbox_labels
// value:map[action:add label:map[id:899202008975144 page_label_name:Intake]
// user:map[id:9979875342055434]]]]
// id:105519898626814 time:1.746604694e+09]] object:page]

//map[entry:[map[changes:[map[field:inbox_labels
// value:map[action:add label:map[id:1023522739876559 page_label_name:ad_id.120227047914070029]
// user:map[id:9392341827487393]]]] id:105519898626814 time:1.746605584e+09]] object:page]

func HandleWebhook(c *fiber.Ctx) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(c.Body(), &raw); err != nil {
		log.Println("Invalid JSON:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Println("Raw Event:", raw)
	var event WebhookDeliveryEvent
	if err := c.BodyParser(&event); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	log.Println("sending message:", event)

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			log.Printf("Received message from %s: %s\n", senderID, msg.Message.Text)
			// if msg.Message.Text != "" || len(msg.Message.Attachments) != 0 {
			// }
			// Save to DB (example)
			err := dbservice.CreateMesseng(&models.Chat{
				SenderID:    senderID,
				UserID:      "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77",
				RecipientID: msg.Recipient.ID,
				JSONMesseng: string(c.BodyRaw()),
			})

			//  query.SetDefault(gormpkg.GetDB())
			//  user := query.Q.User
			// Convert string to int64 first

			var fbID *string
			if senderID != os.Getenv("PAGE_ID") {
				fbID = &senderID
			} else {
				fbID = &msg.Recipient.ID
			}
			// Check if user already exists
			// var existing models.User
			// result := db.Table(models.TableNameUser).Where("facebook_id = ?", fbID).First(&existing)

			// if result.Error != nil && result.Error == gorm.ErrRecordNotFound {

			// User not found, create new one
			gormpkg.GetDB().Table(models.TableNameUser).Create(&customs.UserCustom{
				FacebookID: fbID,
			})
			// }

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "create messeng",
				})

				// Optional reply
				// SendMessage(senderID, "You said: "+msg.Message.Text)
			}

			// Handle postbacks
			if msg.Postback != nil && msg.Postback.Payload != "" {
				log.Printf("Received postback from %s: %s\n", senderID, msg.Postback.Payload)

				// Example response
				SendMessage(senderID, "You clicked: "+msg.Postback.Payload)
			}

			// Handle delivery confirmations
			if msg.Delivery != nil {
				log.Printf("Delivery confirmed for %d message(s): %v\n",
					len(msg.Delivery.Mids), msg.Delivery.Mids)
			}
		}
	}

	// // Handle text messages
	// if msg.Message != nil && msg.Message.Text != "" {
	// 	log.Println("Received message:", senderID, msg.Message.Text)
	// 	dbservice.CreateMesseng(models.Chat{
	// 		SenderID: senderID,
	// 		Message:  msg.Message.Text,
	// 	})
	// 	// Optionally send a reply
	// 	// SendMessage(senderID, "You said: "+msg.Message.Text)
	// }

	// // Handle postback payloads
	// if msg.Postback != nil && msg.Postback.Payload == "SEND_BACK" {
	// 	log.Println("Received postback:", senderID, msg.Postback.Payload)
	// 	SendMessage(senderID, "Thanks for clicking!")
	// }

	return c.SendStatus(fiber.StatusOK)
}
func SendMessage(recipientID, messageText string) error {
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")

	if pageAccessToken == "" {
		return fmt.Errorf(" Missing Token ")
	}
	url := "https://graph.facebook.com/v21.0/me/messages?access_token=" + pageAccessToken

	// Build the message payload
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"messaging_type": "RESPONSE",
		"message": map[string]string{
			"text": messageText,
		},
	}

	// Marshal the payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// // Send the request
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("request failed: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("Facebook API error: %s", resp.Status)
	// }

	return nil
}
