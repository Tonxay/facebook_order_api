package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

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
	var event WebhookEvent
	if err := c.BodyParser(&event); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	log.Println("sending message:", event)

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if msg.Message.Text != "" {
				log.Println("sending message:", senderID, " ", msg.Message.Text)
				// sendMessage(senderID, "You said: "+msg.Message.Text)
			}

			if msg.Postback.Payload == "SEND_BACK" {
				log.Println("sending message:", senderID, " ", msg.Postback.Payload)
				SendMessage(senderID, msg.Message.Text)
			}
		}
	}

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
