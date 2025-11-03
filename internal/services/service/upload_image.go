package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func Uploadimage() {
	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=8252705319817"
	jpegQuality := 1 // Good quality for bill readability
	savePath := "../go-api/images/bills/8252705319817.jpg"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

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

	log.Printf("✅ Successfully saved cropped screenshot to: %s", savePath)
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
