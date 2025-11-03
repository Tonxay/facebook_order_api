package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func Uploadimage() {
	// --- Configuration ---
	// 1. The URL to screenshot
	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=8252705319817"

	// 2. Quality (1-100). Lower numbers = smaller file size (KB)
	jpegQuality := 80

	// 3. Where to save the file
	savePath := "../go-api/images/bills/8252705319817.jpg"
	// ---------------------
	// Set up the browser options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // true = run in background, false = show browser
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// This byte slice will hold the image data
	var buf []byte

	log.Println("Starting screenshot task for:", url)

	// Run the screenshot task
	err := chromedp.Run(ctx,
		chromedp.Emulate(device.IPhone12),          // 1. Emulate iPhone 12
		chromedp.Navigate(url),                     // 2. Go to the URL
		chromedp.Sleep(2*time.Second),              // 3. Wait for page to load
		chromedp.FullScreenshot(&buf, jpegQuality), // 4. Take screenshot
	)

	if err != nil {
		log.Fatalf("Failed during screenshot: %v", err)
	}

	log.Println("Screenshot captured. Saving to:", savePath)

	// --- Call the SaveImage function ---
	if err := SaveImage(buf, savePath); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✅ Successfully saved screenshot to: %s", savePath)
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

package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func Uploadimage() {
	// --- Configuration ---
	// 1. The URL to screenshot
	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=8252705319817"

	// 2. Quality (1-100). Lower numbers = smaller file size (KB)
	jpegQuality := 80

	// 3. Where to save the file
	savePath := "../go-api/images/bills/8252705319817.jpg"
	// ---------------------

	// Set up the browser options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // true = run in background, false = show browser
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// This byte slice will hold the image data
	var buf []byte

	log.Println("Starting screenshot task for:", url)

	// Run the screenshot task
	err := chromedp.Run(ctx,
		chromedp.Emulate(device.IPhone12),          // 1. Emulate iPhone 12
		chromedp.Navigate(url),                     // 2. Go to the URL
		chromedp.Sleep(2*time.Second),              // 3. Wait for page to load
		chromedp.FullScreenshot(&buf, jpegQuality), // 4. Take screenshot
	)

	if err != nil {
		log.Fatalf("Failed during screenshot: %v", err)
	}

	log.Println("Screenshot captured. Saving to:", savePath)

	// --- Call the SaveImage function ---
	if err := SaveImage(buf, savePath); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✅ Successfully saved screenshot to: %s", savePath)
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
