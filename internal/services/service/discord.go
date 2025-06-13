package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Payload structure for Discord webhook
type DiscordWebhook struct {
	Content   string `json:"content"`              // Message content
	Username  string `json:"username,omitempty"`   // Optional: custom username
	AvatarURL string `json:"avatar_url,omitempty"` // Optional: custom avatar
}

func SendDiscordWebhook(webhookURL, message string) error {
	payload := DiscordWebhook{
		Content:  message,
		Username: "Go Bot", // Optional
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned status code: %d", resp.StatusCode)
	}

	fmt.Println("Message sent successfully!")
	return nil
}
