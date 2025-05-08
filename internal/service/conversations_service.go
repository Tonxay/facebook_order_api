package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Struct to unmarshal Facebook conversation data
type Conversation struct {
	ID          string `json:"id"`
	UpdatedTime string `json:"updated_time"`
}

type ConversationsResponse struct {
	Data []Conversation `json:"data"`
}

// GetConversations fetches conversation list from Facebook API
func GetConversations(c *fiber.Ctx) error {
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	pageID := os.Getenv("PAGE_ID")
	apiVersion := "v21.0"

	if pageAccessToken == "" || pageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing page ID or access token",
		})
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/conversations?platform=messenger&access_token=%s", apiVersion, pageID, pageAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to call Facebook API:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch conversations",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	var fbResp ConversationsResponse
	if err := json.Unmarshal(body, &fbResp); err != nil {
		log.Println("Error parsing JSON:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid Facebook API response",
		})
	}

	return c.JSON(fbResp)
}

func GetUserConversation(c *fiber.Ctx) error {
	apiVersion := "v21.0"
	pageID := os.Getenv("PAGE_ID")
	userID := c.Query("user_id")
	platform := c.Query("platform", "messenger")
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")

	if pageID == "" || userID == "" || pageAccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query parameters",
		})
	}

	url := fmt.Sprintf(
		"https://graph.facebook.com/%s/%s/conversations?platform=%s&user_id=%s&access_token=%s",
		apiVersion,
		pageID,
		platform,
		userID,
		pageAccessToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to call Facebook API",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid JSON response",
		})
	}

	return c.JSON(result)
}

func GetMessagesInConversation(c *fiber.Ctx) error {
	apiVersion := "v21.0"
	conversationID := c.Params("conversation_id")
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")

	if conversationID == "" || pageAccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing conversation_id or access token",
		})
	}

	url := fmt.Sprintf(
		"https://graph.facebook.com/%s/%s?fields=messages&access_token=%s",
		apiVersion,
		conversationID,
		pageAccessToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch messages",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	var fbResp map[string]interface{}
	if err := json.Unmarshal(body, &fbResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid JSON response",
		})
	}

	return c.JSON(fbResp)
}

func GetMessageDetails(c *fiber.Ctx) error {
	apiVersion := "v21.0"
	messageID := c.Params("message_id")
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")

	if messageID == "" || pageAccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing message_id or access token",
		})
	}

	url := fmt.Sprintf(
		"https://graph.facebook.com/%s/%s?fields=id,created_time,from,to,message&access_token=%s",
		apiVersion,
		messageID,
		pageAccessToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to request message details",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse JSON",
		})
	}

	return c.JSON(result)
}
