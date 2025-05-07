package main

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

func verifyWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == os.Getenv("VERIFY_TOKEN") {
		return c.SendString(challenge)
	}
	return c.SendStatus(fiber.StatusForbidden)
}

func handleWebhook(c *fiber.Ctx) error {
	var event WebhookEvent
	if err := c.BodyParser(&event); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if msg.Message.Text != "" {
				log.Println("sending message:", senderID, " ", msg.Message.Text)
				// sendMessage(senderID, "You said: "+msg.Message.Text)
			}

			if msg.Postback.Payload == "SEND_BACK" {
				log.Println("sending message:", senderID, " ", msg.Postback.Payload)
				sendMessage(senderID, msg.Message.Text)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
func sendMessage(recipientID, messageText string) error {
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	url := "https://graph.facebook.com/v21.0/me/messages?access_token=" + pageAccessToken

	message := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"messaging_type": "RESPONSE",
		"message": map[string]string{
			"text": messageText,
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return nil
}

// func sendMessage(recipientID, messageText string) {
// 	pageToken := os.Getenv("PAGE_ACCESS_TOKEN")
// 	url := "https://graph.facebook.com/v22.0/me/messages?access_token=" + pageToken

// 	// message := map[string]interface{}{
// 	// 	"recipient": map[string]string{"id": recipientID},
// 	// 	"message":   map[string]string{"text": messageText},
// 	// }
// 	message := map[string]interface{}{
// 		"recipient": map[string]string{
// 			"id": recipientID,
// 		},
// 		"messaging_type": "RESPONSE",
// 		"message": map[string]string{
// 			"text": messageText,
// 		},
// 	}

// 	body, _ := json.Marshal(message)
// 	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Println("Error sending message:", err)
// 	}
// }

func main() {
	app := fiber.New()
	app.Get("/webhook", verifyWebhook)
	app.Post("/webhook", handleWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
