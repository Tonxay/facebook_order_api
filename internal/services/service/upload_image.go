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

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices" // Import consts for UserAgentIPhone
	"github.com/go-rod/rod/lib/proto"
)

func Uploadimage() {
	// log.Println("Starting cron scheduler...")

	// // 1. Set the time zone to Vientiane
	// vientianeZone, err := time.LoadLocation("Asia/Vientiane")
	// if err != nil {
	// 	log.Fatalf("Fatal: Could not load time zone. %v", err)
	// }

	// // 2. Create a new cron scheduler in that time zone
	// c := cron.New(cron.WithLocation(vientianeZone))

	// // 3. Add the Uploadimage function to the schedule
	// //    This string "*/5 * * * *" means "at every 5th minute".
	// //    - "0 * * * *" = "at the start of every hour"
	// //    - "0 9 * * *" = "at 9:00 AM every day"
	// _, err = c.AddFunc("*/1 * * * *", Uploadimage)
	// if err != nil {
	// 	log.Fatalf("Fatal: Could not add cron job. %v", err)
	// }

	// // 4. Start the scheduler
	// c.Start()
	// log.Printf("Scheduler started. Will run job every 5 minutes in %s time.", vientianeZone)

	// TakeScreenshot("8978186958334")

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

// TakeScreenshot navigates to the bill URL, finds the '.bill-content' element,
// captures a JPEG screenshot of only that element, and saves it to a file.
func TakeScreenshot(trackingId string) (string, error) {
	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + trackingId
	saveDir := "../go-api/images/bills/"
	savePath := filepath.Join(saveDir, trackingId+".jpg")

	// 1. Setup a new browser session
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Set up a context with a timeout for the entire operation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a new page and apply settings
	page := browser.MustPage().Context(ctx)

	// --- FIX 1: Set UserAgent *before* navigating ---
	// MustSetUserAgent just takes the user agent string.
	// We do this before navigating so the *first* request has the correct agent.
	page.MustSetUserAgent(devices.IPhoneX.UserAgentEmulation())
	log.Println("Starting screenshot task with go-rod for:", url)

	// --- FIX 1 (cont.): Navigate to the URL ---
	page.MustNavigate(url)
	// We remove MustWaitLoad() because MustNavigate() already waits.
	// --- FIX 2: Use page.Element() instead of WaitElement ---
	// page.Element() is the correct function. It waits for the element
	// to be rendered, finds the first match, and returns it.
	billElement, err := page.Element(".bill-content")
	if err != nil {
		// Try to capture the full page for debugging
		page.MustScreenshot("../go-api/images/bills/debug_error_page.png")
		return "", fmt.Errorf("could not find .bill-content element: %w. Saved debug screenshot", err)
	}

	log.Println("Bill content element found. Capturing screenshot...")

	// 3. Capture the screenshot directly on the element
	// Your code here was already correct.
	// go-rod's Screenshot() when called on an element automatically crops to its bounds.
	buf, err := billElement.Screenshot(
		proto.PageCaptureScreenshotFormatJpeg,
		60, // Quality
	)
	if err != nil {
		return "", fmt.Errorf("failed to capture element screenshot: %w", err)
	}

	log.Println("Screenshot captured and cropped. Saving to:", savePath)

	// 4. Save the image data
	if err := SaveBytes(buf, savePath); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	message := fmt.Sprintf("✅ Successfully saved cropped screenshot to: %s", savePath)
	return message, nil
}

// SaveBytes is a simple helper function to write the byte slice to a file.
func SaveBytes(buf []byte, filePath string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the file
	return os.WriteFile(filePath, buf, 0644)
}

// func scapingImage(trackingId string) string {
// 	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + trackingId
// 	jpegQuality := 60 // Good quality for bill readability
// 	savePath := "../go-api/images/bills/" + trackingId + ".jpg"

// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.Flag("headless", true),
// 		chromedp.Flag("disable-gpu", true),
// 		chromedp.Flag("chromium", true),
// 	)
// 	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// 	defer cancel()

// 	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
// 	defer cancel()

// 	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
// 	defer cancel()

// 	var buf []byte
// 	var billRect map[string]interface{}

// 	log.Println("Starting screenshot task for:", url)

// 	// Get the bill-content's position and dimensions, then take full screenshot
// 	err := chromedp.Run(ctx,
// 		chromedp.Emulate(device.IPhone12),
// 		chromedp.Navigate(url),
// 		chromedp.Sleep(2*time.Second), // Wait for page to load completely

// 		// Get the position and size of the bill-content element
// 		chromedp.Evaluate(`
// 			(() => {
// 				const element = document.querySelector('.bill-content');
// 				if (!element) {
// 					console.error('Could not find .bill-content element');
// 					return null;
// 				}
// 				const rect = element.getBoundingClientRect();
// 				console.log('Bill content rect:', rect);
// 				return {
// 					x: rect.x,
// 					y: rect.y,
// 					width: rect.width,
// 					height: rect.height
// 				};
// 			})()
// 		`, &billRect),

// 		// Take full page screenshot
// 		chromedp.FullScreenshot(&buf, jpegQuality),
// 	)

// 	if err != nil {
// 		log.Fatalf("Failed during screenshot: %v", err)
// 	}

// 	// Check if we found the bill-content element
// 	if billRect == nil {
// 		log.Fatalf("Could not find .bill-content element on the page")
// 	}

// 	x := int(billRect["x"].(float64))
// 	y := int(billRect["y"].(float64))
// 	width := int(billRect["width"].(float64) * 4)
// 	height := int(billRect["height"].(float64) * 3.4)

// 	if width == 0 || height == 0 {
// 		log.Fatalf("Bill content has invalid dimensions: width=%d, height=%d", width, height)
// 	}

// 	log.Printf("Bill content found at: x=%d, y=%d, width=%d, height=%d", x, y, width, height)

// 	// Crop the image to the bill-content
// 	croppedBuf, err := CropImage(buf, x, y, width, height)
// 	if err != nil {
// 		log.Fatalf("Failed to crop image: %v", err)
// 	}

// 	log.Println("Screenshot captured and cropped. Saving to:", savePath)

// 	if err := SaveImage(croppedBuf, savePath); err != nil {
// 		log.Fatalf("Failed to save image: %v", err)
// 	}
// 	message := fmt.Sprintf("✅ Successfully saved cropped screenshot to: %s", savePath)
// 	return message
// 	// sendBillTofaceBook(trackingId)
// }

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
