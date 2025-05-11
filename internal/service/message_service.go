package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	dbservice "go-api/internal/service/db_service"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
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

func HandleWebhook(c *fiber.Ctx) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(c.Body(), &raw); err != nil {
		log.Println("Invalid JSON:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Println("Raw Event:", raw)

	var event WebhookDeliveryEvent
	if err := c.BodyParser(&event); err != nil {
		log.Println("BodyParser error:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Println("Parsed Event:", event)

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			// if msg.Sender == nil || msg.Recipient == nil {
			// 	log.Println("Skipping message with nil sender or recipient")
			// 	continue
			// }

			senderID := msg.Sender.ID
			recipientID := msg.Recipient.ID
			// Store user if not from PAGE_ID
			var fbID string
			if senderID != os.Getenv("PAGE_ID") {
				fbID = senderID
			} else {
				fbID = recipientID
			}

			gormpkg.GetDB().Table(models.TableNameCustomer).Create(&models.Customer{
				FacebookID: fbID,
			})

			if senderID == "" || recipientID == "" {
				log.Println("Skipping message with empty sender or recipient ID")
				continue
			}

			if msg.Message != nil && msg.Message.Text != "" {
				log.Printf("Received message from %s: %s\n", senderID, msg.Message.Text)
			}

			// Save message to DB
			err := dbservice.CreateMesseng(&models.Chat{
				SenderID:    senderID,
				UserID:      "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77", // You might want to dynamically look up user ID
				RecipientID: recipientID,
				JSONMesseng: string(c.BodyRaw()),
			})

			PushToUser(senderID, fiber.Map{
				"customer_id": fbID,
				"message":     c.BodyRaw(),
			})
			if err != nil {
				log.Println("Failed to create message:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "create messeng",
				})
			}

			// Handle postbacks
			if msg.Postback != nil && msg.Postback.Payload != "" {
				log.Printf("Received postback from %s: %s\n", senderID, msg.Postback.Payload)
				SendMessage(senderID, "You clicked: "+msg.Postback.Payload)
			}

			// Handle delivery confirmations
			if msg.Delivery != nil {
				log.Printf("Delivery confirmed for %d message(s): %v\n",
					len(msg.Delivery.Mids), msg.Delivery.Mids)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
