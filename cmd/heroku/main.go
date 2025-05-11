package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// func sendMessage(recipientID, messageText string) error {
// 	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")
// 	url := "https://graph.facebook.com/v21.0/me/messages?access_token=" + pageAccessToken

// 	message := map[string]interface{}{
// 		"recipient": map[string]string{
// 			"id": recipientID,
// 		},
// 		"messaging_type": "RESPONSE",
// 		"message": map[string]string{
// 			"text": messageText,
// 		},
// 	}

// 	jsonData, err := json.Marshal(message)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal message: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return fmt.Errorf("failed to create request: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to send request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("received non-OK response: %s", resp.Status)
// 	}

// 	return nil
// }

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

	// if err := gormpkg.Init("webhook"); err != nil {
	// 	log.Fatalf("❌ Failed to connect to DB: %v", err)
	// }

	// app := fiber.New()

	// // API Routes
	// // api := app.Group(os.Getenv("API_PREFIX"))
	// api.SetupRoutes(app)

	// api.SetupWebsocketRoutes(app)
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }
	// log.Fatal(app.Listen(":" + port))

	app := fiber.New()

	// Facebook webhook verify (GET)
	app.Get("/webhook", func(c *fiber.Ctx) error {
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		verifyToken := os.Getenv("VERIFY_TOKEN") // set in Heroku config

		if mode == "subscribe" && token == verifyToken {
			return c.SendString(challenge)
		}
		return c.SendStatus(fiber.StatusForbidden)
	})

	// Facebook webhook POST
	app.Post("/webhook", func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}
		// Just log or return back for now
		return c.JSON(body)
	})

	// Sample custom API endpoint
	app.Get("/products", func(c *fiber.Ctx) error {
		products := []map[string]interface{}{
			{"id": 1, "name": "Product A"},
			{"id": 2, "name": "Product B"},
		}
		return c.JSON(products)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":3000")

}
