package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"go-api/internal/config/presenters"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SelectTextImage(c *fiber.Ctx) error {
	// Parse multipart form with a file named "image"
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image file is required"})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open image file"})
	}
	defer file.Close()

	// Read file content
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot read image file"})
	}

	// Encode image as base64 string
	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	// Prepare JSON payload for Roboflow API
	payload := map[string]interface{}{
		"api_key": "nfo6qVumn0SgmD6gQP5C", // replace with your Roboflow API key
		"inputs": map[string]interface{}{
			"image": map[string]string{
				"type":  "base64",
				"value": encoded,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to marshal JSON"})
	}

	// Send POST request to Roboflow workflow
	req, err := http.NewRequest("POST", "https://serverless.roboflow.com/infer/workflows/select-text-image/text-v1", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create request"})
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "request to roboflow failed"})
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read roboflow response"})
	}

	// Parse JSON response into a map or struct
	var roboflowResp map[string]interface{}
	if err := json.Unmarshal(respBody, &roboflowResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse roboflow JSON"})
	}

	// Return the Roboflow response directly to client
	return c.JSON(presenters.ResponseSuccess(roboflowResp))

}
