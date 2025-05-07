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

func HandleWebhook(c *fiber.Ctx) error {
	log.Println("sending message:", c.Body())
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

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Facebook API error: %s", resp.Status)
	}

	return nil
}
