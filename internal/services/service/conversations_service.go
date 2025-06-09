package service

import (
	"encoding/json"
	"fmt"
	"go-api/internal/config/middleware"
	custommodel "go-api/internal/pkg/models/custom_model"
	dbservice "go-api/internal/services/db_service"
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

	if pageAccessToken == "" || pageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing page ID or access token",
		})
	}

	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/%s/conversations?fields=participants&access_token=%s",
		pageID, pageAccessToken)

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
	messageID := c.Query("message_id")
	pageId := c.Query("page_id")
	_, token := middleware.CheckPageId(pageId, pageId)
	result, err := GetMessageDetailsFormid(messageID, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	return c.JSON(result)
}

// GetMessageDetailsFormid retrieves message details from the Facebook Graph API
func GetMessageDetailsFormid(messageID, pageAccessToken string) (custommodel.Message, error) {
	apiVersion := "v21.0"
	var result custommodel.Message

	if messageID == "" || pageAccessToken == "" {
		return result, fmt.Errorf("messageID or PAGE_ACCESS_TOKEN is missing")
	}

	// Build URL with messageID and access token
	url := fmt.Sprintf(
		"https://graph.facebook.com/%s/%s?fields=id,created_time,from,to,message&access_token=%s",
		apiVersion,
		messageID,
		pageAccessToken,
	)

	// Send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("failed to send request to Facebook Graph API: %v", err)
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the response body to get more information
		return result, fmt.Errorf("received non-200 response from Facebook API: %d, %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %v", err)
	}

	// Attempt to unmarshal the JSON into the result struct
	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to parse JSON response: %v, response: %s", err, string(body))
	}

	// Return the message details and no error
	return result, nil
}

func GetMessageDetailsPerUser(c *fiber.Ctx) error {
	uerId := c.Params("user_id")
	pageId := os.Getenv("PAGE_ID")
	result, err := dbservice.GetMessengerPerUser(uerId, pageId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	return c.JSON(result)
}
