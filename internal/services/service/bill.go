package service

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Getbill(c *fiber.Ctx) error {
	// 1. Get the tracking number from the URL
	trackingNumber := c.Params("tracking_number")

	// 2. Create the full path to the file
	filePath := filepath.Join("../anousith/images/bills/", trackingNumber+".jpg")

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
func ProxyAnousith(c *fiber.Ctx) error {
	// ດຶງ tracking_number ຈາກ Query String
	trackingNo := c.Query("tracking_number")
	targetURL := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + trackingNo

	req, _ := http.NewRequest("GET", targetURL, nil)
	// ໃສ່ Headers ທີ່ເຮົາກຽມໄວ້
	req.Header.Set("Referer", targetURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; Mobile) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Mobile Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(500).SendString("Fetch Error")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	htmlContent := string(body)

	// ✅ ສັກ <base> tag ເພື່ອໃຫ້ CSS ໂຫຼດໄດ້ (ສຳຄັນຫຼາຍ)
	fixedHTML := strings.Replace(htmlContent, "<head>", "<head><base href=\"https://app.anousith.express/\">", 1)

	c.Set("Access-Control-Allow-Origin", "*") // ໃຫ້ Flutter Web ເຂົ້າເຖິງໄດ້
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(fixedHTML)
}
