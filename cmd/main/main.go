package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go-api/internal/api"
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

	// if err := gormpkg.Init(); err != nil {
	// 	log.Fatalf("❌ Failed to connect to DB: %v", err)
	// }

	app := fiber.New()
	// API Routes
	// api := app.Group(os.Getenv("API_PREFIX"))
	api.SetupRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
