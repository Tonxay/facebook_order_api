package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type WebhookEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Text string `json:"text"`
			} `json:"message,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

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
		log.Println("Failed to parse body:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if msg.Message.Text != "" {
				log.Printf("Received from %s: %s\n", senderID, msg.Message.Text)
				sendMessage(senderID, "Hello")
			}
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func sendMessage(recipientID, messageText string) {
	pageToken := os.Getenv("PAGE_ACCESS_TOKEN")
	url := "https://graph.facebook.com/v17.0/me/messages?access_token=" + pageToken

	message := map[string]interface{}{
		"recipient": map[string]string{"id": recipientID},
		"message":   map[string]string{"text": messageText},
	}

	body, _ := json.Marshal(message)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	app := fiber.New()
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	app.Get("/webhook", verifyWebhook)
	app.Post("/webhook", handleWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(app.Listen(":8080"))
}
