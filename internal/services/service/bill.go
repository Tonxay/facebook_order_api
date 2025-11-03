package service

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func Getbill(c *fiber.Ctx) error {
	// 1. Get the tracking number from the URL
	trackingNumber := c.Params("tracking_number")

	// 2. Create the full path to the file
	filePath := filepath.Join("../go-api/images/bills/", trackingNumber+".jpg")

	// 3. Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Image not found",
		})
	}

	// 4. Send the file. Fiber handles the Content-Type automatically.
	log.Printf("Sending file: %s", filePath)
	return c.SendFile(filePath)
}
