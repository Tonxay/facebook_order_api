package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cons "go-api/internal/config/constant"
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/models/request"
	dbservice "go-api/internal/services/db_service"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices" // Import consts for UserAgentIPhone
	"github.com/go-rod/rod/lib/proto"
	"github.com/gofiber/fiber/v2"
)

func compareLists(list1 []*custommodel.OrderReponseNew, list2 []custommodel.AnousithBill) (matched []custommodel.OrderReponseNew, onlyList1, dataNotMath []string) {
	// normalize to digits only
	normalize := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			if r >= '0' && r <= '9' {
				b.WriteRune(r)
			}
		}
		return b.String()
	}

	// Build map of normalized phone -> present for list1 (no goroutines)
	m1 := make(map[string]bool)
	for _, v := range list1 {
		telStr := fmt.Sprintf("%d", v.Tel)
		n := normalize(telStr)
		if n == "" {
			continue
		}
		m1[n] = true
	}

	// Build map of normalized receiver phone -> trackingId for list2
	m2 := make(map[string]string)
	for _, v := range list2 {
		n := normalize(v.ReceiverPhone)
		if n == "" {
			continue
		}
		m2[n] = v.TrackingId
	}

	// Check list1 and assign tracking IDs when matched
	for _, v := range list1 {
		telStr := fmt.Sprintf("%d", v.Tel)
		normTel := normalize(telStr)

		found := false
		for recvNorm, tracking := range m2 {
			if normTel == "" {
				continue
			}
			if strings.HasSuffix(recvNorm, normTel) || strings.HasSuffix(normTel, recvNorm) {
				for _, item := range list2 {

					if item.TrackingId == tracking {
						v.AnousithBillData = item
					}

				}

				v.TrackingNumber = tracking
				v.LikeTackingNumberURL = "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + tracking
				matched = append(matched, *v)
				found = true
				break
			}
		}

		if !found {
			onlyList1 = append(onlyList1, telStr)
		}
	}

	// Check list2 for items not in list1
	for _, v := range list2 {
		n := normalize(v.ReceiverPhone)
		if n == "" {
			dataNotMath = append(dataNotMath, v.ReceiverPhone)
			continue
		}
		if !m1[n] {
			dataNotMath = append(dataNotMath, v.ReceiverPhone)
		}
	}

	return
}

// "ordered", // ສັ່ງຊື້ (Order placed)
// "waiting_to_pack", // ລໍຖ້າແພັກເຄື່ອງ (Waiting to pack)
// "packed", // ແພັກເຄື່ອງແລ້ວ (Packed)
// "shipped", // ສົ່ງແລ້ວ (Shipped)
// "customer_bill_notified", // ແຈ້ງບິນໃຫ້ລູກຄ້າແລ້ວ (Customer bill notified)
// "delivery_complete", // ສົ່ງສຳເລັດ (Delivery complete)
// "payment_completed", // ສຳລະເງິນແລ້ວ (Payment completed)
// "order_cancelled", // ຍົກເລີກອໍເດີ (Order cancelled)
// "return_to_sender", // ພັດສະດຸຕີກັບ (Returned to sender)
// "customer_notified", // ແຈ້ງລູກແລ້ວ (Customer notified)

func GetOrderbillInAnousith(c *fiber.Ctx) error {
	token := AnousithLoging("92339355", "s0987654")
	order, _ := dbservice.GetOrders(gormpkg.GetDB(), request.StatusOrderRequest{
		IsCancel:   false,
		ShippingID: "3531696c-af5b-4d73-b677-dbd9e1bb4f1b",
		Statuses:   []string{"shipped"},
	})

	if len(order) == 0 {
		return fmt.Errorf("No orders found")
	}
	amountBills := AnousithOrder(token)

	matched, onlyL1, dataNotMath := compareLists(order, amountBills.Data.ItemsV2.Data)

	fmt.Println("Matched:", matched)
	fmt.Println("Only in List1:", onlyL1)
	fmt.Println("Only in data_not_matched:", dataNotMath)
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(fiber.Map{
		"count_orders":                         len(order),
		"count_order_matched":                  len(matched),
		"count_order_not_matched":              len(onlyL1),
		"count_anousith_bill_data_not_matched": len(dataNotMath),
		"matched":                              matched,
		"order_not_matched":                    onlyL1,
		"anousith_bill_data_not_matched":       dataNotMath,
	}))

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

func AnousithLoging(tel string, password string) string {

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

	// Parse JSON response to extract accessToken safely
	var respJSON map[string]interface{}
	if err := json.Unmarshal(body, &respJSON); err != nil {
		log.Fatal("Error parsing response JSON:", err)
	}

	accessToken := ""
	if data, ok := respJSON["data"].(map[string]interface{}); ok {
		if cl, ok := data["customerLogin"].(map[string]interface{}); ok {
			if at, ok := cl["accessToken"].(string); ok {
				accessToken = at
			}
		}
	}

	if accessToken == "" {
		log.Printf("accessToken not found in response: %s", string(body))
	}

	return accessToken

}

func AnousithOrder(token string) custommodel.ItemsV2Response {

	payload := map[string]interface{}{
		"operationName": "ItemsV2",
		"variables": map[string]interface{}{
			"where": map[string]interface{}{
				"multipleItemStatus": []string{
					"TRANSIT_TO_DEST_BRANCH", // ກຳລັງຂົນສົ່ງໄປປາຍທາງ
					// "TRANSIT_TO_ORIGIN_BRANCH", // ເຄື່ອງກຳລັງຕີກັບ
					"DEST_BRANCH_RECEIVED_FORWARD",    // ສາຂາປາຍທາງຮັບເຄື່ອງແລ້ວ
					"ORIGIN_BRANCH_RECEIVED_BACKWARD", // ສາຂາຕົ້ນທາງຮັບເຄື່ອງແລ້ວ
					"DEST_BRANCH_RECEIVED_BACKWARD",   // ສາຂາປາຍທາງຮັບເຄື່ອງແລ້ວ
					"ORIGIN_BRANCH_RECEIVED_FORWARD",  // ສາຂາຕົ້ນທາງຮັບເຄື່ອງແລ້ວ
					"COMPLETED",                       // ສົ່ງສຳເລັດແລ້ວ
				},
				"originReceiveDate_gte": func() string {
					now := time.Now().UTC()
					start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
					return start.Format("2006-01-02")
				}(),
				"originReceiveDate_lt": func() string {
					now := time.Now().UTC()
					start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
					end := start.AddDate(0, 1, 0) // first day of next month (exclusive)
					return end.Format("2006-01-02")
				}(),
				"searchMultipleCOD": []string{"0", "1"},
				"customerId":        6216826,
				"isDeleted":         0,
			},
			"orderBy": "originReceiveDate_DESC",
			"skip":    0,
			"limit":   1000,
		},
		"query": `query ItemsV2($where: ItemV2WhereInput, $skip: Int, $noLimit: Boolean, $limit: Int, $orderBy: OrderByItem) {
  itemsV2(
    where: $where
    skip: $skip
    noLimit: $noLimit
    limit: $limit
    orderBy: $orderBy
  ) {
    total
    data {
      _id
      trackingPlatform
      trackingId
      itemName
      itemValueKIP
      itemValueTHB
      itemValueUSD
      realItemValueKIP
      realItemValueTHB
      realItemValueUSD
      receiverName
      receiverPhone
      description
      isSummary
      destSendDate
      charge_on_shop
      itemStatus
      contactStatus
      originSendDate
      width
      weight
      isCod
      isExtraItem
      packagePrice
      originReceiveDate
      destReceiveDate
      sendCompleteDate
      isBackward
      billNumber
      originProvinceId {
        provinceName
      }
      destProvinceId {
        provinceName
      }
      originBranchId {
        branch_name
      }
      destBranchId {
        branch_name
        branch_address
        contactInfo
      }
      customerId {
        id_list
        full_name
        contact_info
      }
      createdBy {
        first_name
        phone_number
      }
      originReceiveBy {
        first_name
        phone_number
      }
      providedBy {
        _id
      }
    }
  }
}`,
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

	req.Header.Set("Authorization", token)
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
	log.Println("Response body:", string(body))
	if err != nil {
		log.Fatal(err)
	}
	var response custommodel.ItemsV2Response
	json.Unmarshal(body, &response)

	return response
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

func sendBillTofaceBook(order OrderQueryParams, token string) error {

	facebookUrl := "https://graph.facebook.com/v18.0/me/messages?access_token=" + token

	// msg := map[string]any{
	// 	"recipient": map[string]any{
	// 		"id": "23933715116288753",
	// 	},
	// 	"message": map[string]any{
	// 		"attachment": map[string]any{
	// 			"type": "image",
	// 			"payload": map[string]any{
	// 				"url":         "https://api.chat-dd.uk/bill/" + trackingId,
	// 				"is_reusable": true,
	// 			},
	// 		},
	// 	},
	// 	"messaging_type": "MESSAGE_TAG",
	// 	"tag":            "POST_PURCHASE_UPDATE",
	// }

	textMessage := "ສະບາຍດີລູກຄ້າ: ຝາກເຄື່ອງລົງ ອານຸສິດ ໃຫ້ເເລ້ວເຈົ້າ \n" +
		"ເລກຕິດຕາມ: " + order.TrackingNumber + "\n" +
		"ຊື່: " + order.CustomerName + "\n" +
		"ເບີ: " + order.CustomerTel + "\n" +
		"ເບີສາຂາ: " + order.ContactInfo + "\n" +
		"ຝາກລົງ: " + order.BranchAddress

	textMessageThank := "ຂອບໃຈຫຼາຍໆ🙏🏻💖ສຳລັບການອຸດໜູນ! \n - ໃຫ້ລູກຄ້າກົດເບີ່ງຮູບ ເພື່ອບັນທຶກເອົາບິນ \n - ຫຼັງຈາກໄດ້ຮັບສິນຄ້າ ລົບກວນຖ່າຍວີດິໂອຕອນແກະເຄື່ອງໄວ້ແນ່ເດີ \n - ໃນກໍລະນີເຄື່ອງມີບັນຫາທາງຮ້ານເຮົາຈຶ່ງຈະຮັບຜິດຊອບໃຫ້ເດີເຈົ້າ. 🤗"

	msgText := map[string]any{
		"recipient": map[string]any{
			"id": order.CustomerID,
		},
		"message": map[string]any{
			"text": textMessage, // ✅ Changed from 'attachment' to 'text'
		},
		"messaging_type": "MESSAGE_TAG",
		"tag":            "POST_PURCHASE_UPDATE",
	}
	msgTextThank := map[string]any{
		"recipient": map[string]any{
			"id": order.CustomerID,
		},
		"message": map[string]any{
			"text": textMessageThank, // ✅ Changed from 'attachment' to 'text'
		},
		"messaging_type": "MESSAGE_TAG",
		"tag":            "POST_PURCHASE_UPDATE",
	}

	msg := map[string]any{
		"recipient": map[string]any{
			"id": order.CustomerID,
		},
		"message": map[string]any{
			"attachment": map[string]any{
				"type": "template",
				"payload": map[string]any{
					"template_type": "generic",
					"elements": []map[string]any{
						{
							"title":     "ໃບບິນຝາກພັດສະດຸ (ອານຸສິດ)",
							"image_url": "https://scontent.fvte1-1.fna.fbcdn.net/v/t1.6435-9/76994883_478719836071281_7632067609801785344_n.jpg?_nc_cat=102&ccb=1-7&_nc_sid=a5f93a&_nc_eui2=AeEMYTkT3fPflhxKDc1ovMS3SfthehKY9-xJ-2F6Epj37MjKYxXpDwqyzg8fSIqnekInSknoha-loEr-U-D08w78&_nc_ohc=ESG-LZEfBPEQ7kNvwEZaJIw&_nc_oc=Adm1npqBFHO_iyW6lpJczxoTlh36U2R--uH6bC1osVskDmrFe2lTpUi3tkiSA-dLP10&_nc_zt=23&_nc_ht=scontent.fvte1-1.fna&_nc_gid=LM_fVkkyNUTTcfJrl8Tuog&oh=00_AfrsIibwGtSWLvQOI6Q5Li4wEKUpqvLmomKpVFyg3NQSBQ&oe=698332AA",
							"subtitle":  "ເລກຕິດຕາມ: " + order.TrackingNumber + "\nຊື່: " + order.CustomerName + "ເບີ: " + order.CustomerTel,
							"buttons": []map[string]any{
								{
									"type":                 "web_url",
									"url":                  "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + order.TrackingNumber,
									"title":                "ກົດເບິ່ງຮູບໃບບິນ",
									"webview_height_ratio": "full",
								},
							},
						},
					},
				},
			},
		},
		// ✅ ຕ້ອງຢູ່ບ່ອນນີ້! (Root level)
		"messaging_type": "MESSAGE_TAG",
		"tag":            "POST_PURCHASE_UPDATE",
	}
	var err error
	err = send(msgText, facebookUrl)
	err = send(msg, facebookUrl)
	err = send(msgTextThank, facebookUrl)

	if err != nil {
		return err
	}
	db := gormpkg.GetDB()
	status := "customer_bill_notified"
	oldStatus, ok := cons.OrderStatusTransitions[status]
	if !ok {
		return fiber.NewError(400, "not found status")
	}
	_, err = dbservice.UpdateStatusOrder(db, order.OrderID, status, oldStatus, "", order.UserID, order.TrackingNumber)
	if err != nil {
		return err
	}
	return nil
}

func send(msg map[string]any, facebookUrl string) error {
	// 2. Marshal it to JSON (same as before)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, facebookUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	// 3. Send the request using the default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 4. Read the response (same as before)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	log.Println("✅ Successfully sended:", req)
	return nil
}

// func sendBillTofaceBook(trackingId string, token string, customerId string) {

// 	facebookUrl := "https://graph.facebook.com/v18.0/me/messages?access_token=EAARDcwZBMbeQBPZC3KIzAjxMn2HOtRv498MYVpKxSc16Du03srxwRC8M26b9CLMY4qZCvVoj1e11HgZCNWEDCKCviosA831C4hn7IHPUUUfrPvUFKUjLBA4f9bzpRswtvg9RtveJT6ZCB6nI7FZA7wWFFxT0K0a6A8ft4pdsxW0sm1a2AkCyiz2lb1kyYPl1teU1QZD"

// 	// msg := map[string]any{
// 	// 	"recipient": map[string]any{
// 	// 		"id": "9608668882534684",
// 	// 	},
// 	// 	"message": map[string]any{
// 	// 		"attachment": map[string]any{
// 	// 			"type": "image",
// 	// 			"payload": map[string]any{
// 	// 				"url":         "https://api.chat-dd.uk/bill/" + trackingId,
// 	// 				"is_reusable": true,
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	"messaging_type": "MESSAGE_TAG",
// 	// 	"tag":            "POST_PURCHASE_UPDATE",
// 	// }

// 	// msg := map[string]any{
// 	// 	"recipient": map[string]any{
// 	// 		"id": "23933715116288753",
// 	// 	},
// 	// 	"message": map[string]any{
// 	// 		"attachment": map[string]any{
// 	// 			"type": "template",
// 	// 			"payload": map[string]any{
// 	// 				"template_type": "generic",
// 	// 				"elements": []map[string]any{
// 	// 					{
// 	// 						"title":     "ໃບບິນຝາກພັດສະດຸ (ອານຸສິດ)",
// 	// 						"image_url": "https://app.anousith.express/static/media/logo-app.08726185419ef9a8e073.png", // ຮູບໂລໂກ້ ຫຼື ໄອຄອນ
// 	// 						"subtitle":  "ເລກຕິດຕາມ: " + trackingId,
// 	// 						"buttons": []map[string]any{
// 	// 							{
// 	// 								"type":  "web_url",
// 	// 								"url":   "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + trackingId,
// 	// 								"title": "ກົດເບິ່ງໃບບິນຕົວຈິງ",
// 	// 								// ✅ ເຮັດໃຫ້ເປີດຂຶ້ນມາເຕັມຈໍພາຍໃນ Messenger
// 	// 								"webview_height_ratio": "full",
// 	// 								"messenger_extensions": true,
// 	// 							},
// 	// 						},
// 	// 					},
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 	},
// 	// }

// 	// 2. Marshal it to JSON (same as before)
// 	jsonData, err := json.Marshal(msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req, err := http.NewRequest(http.MethodPost, facebookUrl, bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// 3. Send the request using the default client
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	// 4. Read the response (same as before)
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Status:", resp.Status)
// 	fmt.Println("Response Body:", string(body))
// 	log.Println("✅ Successfully sended:", req)
// }

// func ScapingImage(c *fiber.Ctx) error {
// 	var tracking_number = c.Query("tracking_number", "")
// 	if tracking_number == "" {
// 		return c.Status(fiber.ErrBadRequest.Code).SendString("❌ Missing 'tracking_number'")
// 	}

// 	url := "https://app.anousith.express/landing/search_tracking/bill_share?tacking_number=" + tracking_number
// 	savePath := "../anousith/images/bills/" + tracking_number + ".jpg"

// 	// 1. ປັບປຸງ Allocator Options
// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.NoSandbox,
// 		chromedp.DisableGPU,
// 		chromedp.Flag("disable-dev-shm-usage", true),
// 	)
// 	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// 	defer cancel()

// 	ctx, cancel := chromedp.NewContext(allocCtx)
// 	defer cancel()

// 	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
// 	defer cancel()

// 	var buf []byte
// 	var billRect map[string]interface{}

// 	// 2. ໃຊ້ Emulate ດ້ວຍ Scale ທີ່ສູງຂຶ້ນເພື່ອຄວາມຊັດ (Scale 3.0 = 3x resolution)
// 	highResDevice := device.Info{
// 		Name:      "HighResDevice",
// 		UserAgent: "Mozilla/5.0...",
// 		Width:     500,
// 		Height:    1000,
// 		Scale:     3.0, // ເພີ່ມຈາກ 1.0 ເປັນ 100 ຈະຊັດຂຶ້ນຫຼາຍ
// 		Landscape: false,
// 		Mobile:    true,
// 	}

// 	err := chromedp.Run(ctx,
// 		chromedp.Emulate(highResDevice),
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(`.bill-content`, chromedp.ByQuery), // ລໍຖ້າໃຫ້ Element ຂຶ້ນມາແທນການ Sleep
// 		chromedp.Sleep(1*time.Second),                           // ຖ້າມີ QR code ໃຫ້ລໍຖ້າແປັບໜຶ່ງ

// 		chromedp.Evaluate(`
//             (() => {
//                 const element = document.querySelector('.bill-content');
//                 if (!element) return null;
//                 const rect = element.getBoundingClientRect();
//                 const dpr = window.devicePixelRatio || 1;
//                 return {
//                     x: rect.x * dpr,
//                     y: rect.y * dpr,
//                     width: rect.width * dpr,
//                     height: rect.height * dpr
//                 };
//             })()
//         `, &billRect),

// 		// ປັບ Quality ເປັນ 100 (ສູງສຸດ)
// 		chromedp.FullScreenshot(&buf, 100),
// 	)

// 	if err != nil || billRect == nil {
// 		return c.Status(500).SendString("Capture Failed")
// 	}

// 	// ຄິດໄລ່ຕຳແໜ່ງຕາມ Scale ທີ່ໄດ້ມາຈາກ JS ໂດຍກົງ
// 	x := int(billRect["x"].(float64))
// 	y := int(billRect["y"].(float64))
// 	width := int(billRect["width"].(float64))
// 	height := int(billRect["height"].(float64))

// 	// 3. Crop ດ້ວຍຄວາມລະອຽດສູງ
// 	croppedBuf, err := CropImage(buf, x, y, width, height)
// 	if err != nil {
// 		return c.Status(500).SendString("Crop Failed")
// 	}

// 	if err := SaveImage(croppedBuf, savePath); err != nil {
// 		return c.Status(500).SendString("Save Failed")
// 	}

//		sendBillTofaceBook(tracking_number)
//		return c.SendString("✅ Saved High-Res Image")
//	}

var (
	GlobalAllocCtx context.Context
	Sem            = make(chan struct{}, 5) // ຄຸມຄິວ 5 Tasks
)

func InitBrowser() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("headless", true),
	)
	// ກຳນົດຄ່າໃຫ້ຕົວແປ Global
	GlobalAllocCtx, _ = chromedp.NewExecAllocator(context.Background(), opts...)
}

type OrderQueryParams struct {
	CustomerID     string `query:"customer_id" validate:"required"`
	PageID         string `query:"page_id" validate:"required"`
	OrderID        string `query:"order_id" validate:"required"`
	UserID         string
	TrackingNumber string  `query:"tracking_number" validate:"required"`
	PackagePrice   float64 `query:"package_price" validate:"required"`
	CustomerName   string  `query:"customer_name" validate:"required"`
	CustomerTel    string  `query:"customer_tel" validate:"required"`
	FreeShipping   bool    `query:"free_shipping" validate:"required"`
	COD            bool    `query:"cod" validate:"required"`
	Platform       string  `query:"platform" validate:"required"`
	BranchAddress  string  `query:"branch_address" validate:"required"`
	ContactInfo    string  `query:"contactinfo" validate:"required"`
}

func ScapingImage(c *fiber.Ctx) error {
	// log.Println("Starting ScapingImage handler")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}

	var params OrderQueryParams
	err := c.QueryParser(&params)

	if err != nil {
		return fiber.NewError(400, " failed get data error")
	}

	params.UserID = userID

	trackingNo := params.TrackingNumber
	if trackingNo == "" {
		return fiber.NewError(400, "❌ Missing tracking_number")
	}
	if params.CustomerID == "" {
		return fiber.NewError(400, "❌ Missing customer_id")
	}

	if params.Platform == "whatsapp" {
		return fiber.NewError(402, "❌ Missing whatapp")
	}
	_, token := middleware.CheckPageId(params.PageID, params.PageID)

	sendBillTofaceBook(params, token) // ເອີ້ນໃຊ້ຈາກ utils.go

	return c.Status(200).SendString("✅ Success")
}

// // --- Helper Functions ---

// func CropImage(input []byte, x, y, w, h int) ([]byte, error) {
// 	src, _, err := image.Decode(bytes.NewReader(input))
// 	if err != nil {
// 		return nil, err
// 	}

// 	if sub, ok := src.(interface {
// 		SubImage(r image.Rectangle) image.Image
// 	}); ok {
// 		img := sub.SubImage(image.Rect(x, y, x+w, y+h))
// 		buf := new(bytes.Buffer)
// 		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
// 		return buf.Bytes(), err
// 	}
// 	return nil, fmt.Errorf("crop not supported")
// }

// func SaveImage(data []byte, path string) error {
// 	// ສ້າງ Folder ຖ້າຍັງບໍ່ມີ
// 	dir := "/anousith/images/bills" // ປ່ຽນຕາມ Path ໃນ Container
// 	os.MkdirAll(dir, 0755)
// 	return os.WriteFile(path, data, 0644)
// }
