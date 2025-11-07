package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func Uploadimage() {
	scapingImage("8978186958334")
}

// SaveImage saves a data slice (like an image) to the specified file path.
// It automatically creates all necessary parent directories if they do not exist.
func SaveImage(data []byte, filePath string) error {

	// 1. Get the directory part of the path
	dir := filepath.Dir(filePath)

	// 2. Create the directory (and all parents) if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 3. Write the file to the specified path
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// CropImage crops an image to the specified rectangle
func CropImage(imgData []byte, x, y, width, height int) ([]byte, error) {
	// Decode the image
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	log.Printf("Original image format: %s, bounds: %v", format, img.Bounds())

	// Ensure crop coordinates are within image bounds
	bounds := img.Bounds()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+width > bounds.Dx() {
		width = bounds.Dx() - x
	}
	if y+height > bounds.Dy() {
		height = bounds.Dy() - y
	}

	// Define crop rectangle
	rect := image.Rect(x, y, x+width, y+height)

	log.Printf("Cropping to: %v", rect)

	// Create a new image with the cropped region
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	croppedImg := img.(subImager).SubImage(rect)

	// Encode back to JPEG with good quality for text readability
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, croppedImg, &jpeg.Options{Quality: 1}); err != nil {
		return nil, fmt.Errorf("failed to encode cropped image: %w", err)
	}

	log.Printf("Cropped image size: %d bytes", buf.Len())

	return buf.Bytes(), nil
}

func anousithLoging(tel string, password string) {

	payload := map[string]interface{}{
		"operationName": "CustomerLogin",
		"variables": map[string]interface{}{
			"where": map[string]interface{}{
				"username": tel,      // Use the parameter here
				"password": password, // And here
			},
		},
		"query": "mutation CustomerLogin($where: CustomerLoginInput!) {\n  customerLogin(where: $where) {\n    accessToken\n    data {\n      id_list\n      full_name\n      profile_img\n      status\n      contact_info\n      address\n      village\n      district {\n        id_list\n        title\n      }\n      state {\n        provinceName\n        id_state\n      }\n      Bank_KIP\n      BANK_THB\n      BANK_USD\n      BANK_NAME\n      gender\n      isActive\n      isVerify\n    }\n  }\n}",
	}

	// 2. "Marshal" the map into a JSON byte slice
	// This will safely handle any special characters in your params.
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error creating JSON payload:", err)
	}

	// 'jsonData' is now ready to be used in an HTTP request
	fmt.Println("--- Generated JSON Byte Slice ---")
	// We use bytes.Buffer to pretty-print the JSON for demonstration
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonData, "", "  ")

	fmt.Println(prettyJSON.String()) // (Your JSON data)
	apiUrl := "https://pro.api.anousith.express/graphql"

	// 1. Create a new request, but don't send it yet
	// (Note: http.MethodPost is just a constant for "POST")
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set your headers (THE IMPORTANT PART)
	// The "application/json" from your old code *is* a header, so set it here.
	req.Header.Set("Content-Type", "application/json")

	// Add any other headers you need
	req.Header.Set("referer", "https://app.anousith.express/")

	// 3. Send the request using the default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 4. Read the response (same as before)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

}

func scapingImage(trackingId string) string {
	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + trackingId
	jpegQuality := 60 // Good quality for bill readability
	savePath := "../go-api/images/bills/" + trackingId + ".jpg"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("chromium", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()

	ctx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	// defer cancel()

	ctx, _ = context.WithTimeout(ctx, 120*time.Second)
	// defer cancel()

	var buf []byte
	var billRect map[string]interface{}

	log.Println("Starting screenshot task for:", url)

	// Get the bill-content's position and dimensions, then take full screenshot
	err := chromedp.Run(ctx,
		chromedp.Emulate(device.IPhone12),
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Wait for page to load completely

		// Get the position and size of the bill-content element
		chromedp.Evaluate(`
			(() => {
				const element = document.querySelector('.bill-content');
				if (!element) {
					console.error('Could not find .bill-content element');
					return null;
				}
				const rect = element.getBoundingClientRect();
				console.log('Bill content rect:', rect);
				return {
					x: rect.x,
					y: rect.y,
					width: rect.width,
					height: rect.height
				};
			})()
		`, &billRect),

		// Take full page screenshot
		chromedp.FullScreenshot(&buf, jpegQuality),
	)

	if err != nil {
		log.Fatalf("Failed during screenshot: %v", err)
	}

	// Check if we found the bill-content element
	if billRect == nil {
		log.Fatalf("Could not find .bill-content element on the page")
	}

	x := int(billRect["x"].(float64))
	y := int(billRect["y"].(float64))
	width := int(billRect["width"].(float64) * 4)
	height := int(billRect["height"].(float64) * 3.4)

	if width == 0 || height == 0 {
		log.Fatalf("Bill content has invalid dimensions: width=%d, height=%d", width, height)
	}

	log.Printf("Bill content found at: x=%d, y=%d, width=%d, height=%d", x, y, width, height)

	// Crop the image to the bill-content
	croppedBuf, err := CropImage(buf, x, y, width, height)
	if err != nil {
		log.Fatalf("Failed to crop image: %v", err)
	}

	log.Println("Screenshot captured and cropped. Saving to:", savePath)

	if err := SaveImage(croppedBuf, savePath); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	message := fmt.Sprintf("✅ Successfully saved cropped screenshot to: %s", savePath)
	return message
	// sendBillTofaceBook(trackingId)
}

func sendBillTofaceBook(trackingId string) {

	facebookUrl := "https://graph.facebook.com/v18.0/me/messages?access_token=EAARDcwZBMbeQBPZC3KIzAjxMn2HOtRv498MYVpKxSc16Du03srxwRC8M26b9CLMY4qZCvVoj1e11HgZCNWEDCKCviosA831C4hn7IHPUUUfrPvUFKUjLBA4f9bzpRswtvg9RtveJT6ZCB6nI7FZA7wWFFxT0K0a6A8ft4pdsxW0sm1a2AkCyiz2lb1kyYPl1teU1QZD"

	msg := map[string]any{
		"recipient": map[string]any{
			"id": "23933715116288753",
		},
		"message": map[string]any{
			"attachment": map[string]any{
				"type": "image",
				"payload": map[string]any{
					"url":         "https://api.chat-dd.uk/bill/" + trackingId,
					"is_reusable": true,
				},
			},
		},
		"messaging_type": "MESSAGE_TAG",
		"tag":            "POST_PURCHASE_UPDATE",
	}

	// 2. Marshal it to JSON (same as before)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, facebookUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	// 3. Send the request using the default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 4. Read the response (same as before)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	log.Println("✅ Successfully sended:", req)
}
