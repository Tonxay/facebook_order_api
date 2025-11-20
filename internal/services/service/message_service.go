package service

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

// --- 1. Defince Structs (ประกาศโครงสร้างข้อมูลให้ตรงกับ Facebook JSON) ---
type FacebookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        string      `json:"id"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"` // Pointer เพื่อเช็ค nil ได้
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID      string    `json:"mid"`
	Text     string    `json:"text"`
	Referral *Referral `json:"referral,omitempty"` // จุดสำคัญ: ข้อมูล Ad อยู่ที่นี่
}

type Referral struct {
	Source string `json:"source"`
	Type   string `json:"type"`
	AdID   string `json:"ad_id"` // <--- ID ของโฆษณาที่เราต้องการ

	// [เพิ่ม]: รับข้อมูล Ads Context Data
	AdsContextData *AdsContextData `json:"ads_context_data,omitempty"`
}

// [เพิ่ม]: Struct ใหม่สำหรับรับ ad_title
type AdsContextData struct {
	AdTitle  string `json:"ad_title"`
	PostID   string `json:"post_id"`   // เก็บ Post ID ด้วยก็ได้เผื่อใช้
	VideoURL string `json:"video_url"` // เก็บ Video URL ด้วยก็ได้
}

// type WebhookDeliveryEvent struct {
// 	Object string `json:"object"`
// 	Entry  []struct {
// 		ID        string `json:"id"`
// 		Time      int64  `json:"time"`
// 		Messaging []struct {
// 			Sender struct {
// 				ID string `json:"id"`
// 			} `json:"sender"`
// 			Recipient struct {
// 				ID string `json:"id"`
// 			} `json:"recipient"`

// 			Timestamp float64 `json:"timestamp"`

// 			Message struct {
// 				Mid         string `json:"mid"`
// 				Text        string `json:"text,omitempty"`
// 				Attachments []struct {
// 					Type    string `json:"type"`
// 					Payload struct {
// 						URL string `json:"url"`
// 					} `json:"payload"`
// 				} `json:"attachments,omitempty"`
// 			} `json:"message,omitempty"`

// 			Delivery *struct {
// 				Mids      []string `json:"mids"`
// 				Watermark float64  `json:"watermark"`
// 			} `json:"delivery,omitempty"`

// 			Postback *struct {
// 				Payload string `json:"payload"`
// 			} `json:"postback,omitempty"`
// 		} `json:"messaging"`
// 	} `json:"entry"`
// }

// func VerifyWebhook(c *fiber.Ctx) error {
// 	mode := c.Query("hub.mode")
// 	token := c.Query("hub.verify_token")
// 	challenge := c.Query("hub.challenge")

// 	if mode == "subscribe" && token == os.Getenv("VERIFY_TOKEN") {
// 		return c.SendString(challenge)
// 	}
// 	return c.SendStatus(fiber.StatusForbidden)
// }

// func SendMessage(recipientID, messageText string) error {
// 	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")

// 	if pageAccessToken == "" {
// 		return fmt.Errorf(" Missing Token ")
// 	}
// 	url := "https://graph.facebook.com/v21.0/me/messages?access_token=" + pageAccessToken

// 	// Build the message payload
// 	payload := map[string]interface{}{
// 		"recipient": map[string]string{
// 			"id": recipientID,
// 		},
// 		"messaging_type": "RESPONSE",
// 		"message": map[string]string{
// 			"text": messageText,
// 		},
// 	}

// 	// Marshal the payload to JSON
// 	body, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("json marshal error: %w", err)
// 	}

// 	// Prepare the HTTP request
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

// 	if err != nil {
// 		return fmt.Errorf("request creation error: %w", err)
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")

// 	// // Send the request
// 	// client := &http.Client{}
// 	// resp, err := client.Do(req)
// 	// if err != nil {
// 	// 	return fmt.Errorf("request failed: %w", err)
// 	// }
// 	// defer resp.Body.Close()

// 	// if resp.StatusCode != http.StatusOK {
// 	// 	return fmt.Errorf("Facebook API error: %s", resp.Status)
// 	// }

// 	return nil
// }

// func HandleWebhook(c *fiber.Ctx) error {
// 	// var raw map[string]interface{}
// 	// if err := json.Unmarshal(c.Body(), &raw); err != nil {
// 	// 	log.Println("Invalid JSON:", err)
// 	// 	return c.SendStatus(fiber.StatusBadRequest)
// 	// }

// 	// log.Println("Raw Event:", raw)

// 	// var event WebhookDeliveryEvent
// 	// if err := c.BodyParser(&event); err != nil {
// 	// 	log.Println("BodyParser error:", err)
// 	// 	return c.SendStatus(fiber.StatusBadRequest)
// 	// }

// 	// log.Println("Parsed Event:", event)
// 	// var user *models.Customer

// 	// for _, entry := range event.Entry {
// 	// 	for _, msg := range entry.Messaging {
// 	// 		// if msg.Sender == nil || msg.Recipient == nil {
// 	// 		// 	log.Println("Skipping message with nil sender or recipient")
// 	// 		// 	continue
// 	// 		// }
// 	// 		// if msg.Message.Text != "" || len(msg.Message.Attachments) >= 1 {

// 	// 		// 	senderID := msg.Sender.ID
// 	// 		// 	recipientID := msg.Recipient.ID
// 	// 		// 	// Store user if not from PAGE_ID
// 	// 		// 	var fbID string

// 	// 		// 	pageId, token := middleware.CheckPageId(senderID, recipientID)
// 	// 		// 	if senderID != pageId {
// 	// 		// 		fbID = senderID
// 	// 		// 	} else {
// 	// 		// 		fbID = recipientID
// 	// 		// 	}

// 	// 		// 	err1 := gormpkg.GetDB().Table(models.TableNameCustomer).Create(&models.Customer{
// 	// 		// 		FacebookID:  fbID,
// 	// 		// 		PageID:      pageId,
// 	// 		// 		PhoneNumber: 0,
// 	// 		// 	}).Error

// 	// 		// 	gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).UpdateColumns(&models.Customer{
// 	// 		// 		UpdatedAt: time.Now(),
// 	// 		// 	})

// 	// 		// 	if err1 == nil {
// 	// 		// 		var fullnam string
// 	// 		// 		message, _ := GetMessageDetailsFormid(msg.Message.Mid, token)
// 	// 		// 		fmt.Println(message)
// 	// 		// 		if message.From.ID == fbID {
// 	// 		// 			fullnam = message.From.Name
// 	// 		// 		} else {
// 	// 		// 			fullnam = message.To.Data[0].Name
// 	// 		// 		}
// 	// 		// 		gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).UpdateColumns(&models.Customer{
// 	// 		// 			FirstName: fullnam,
// 	// 		// 			UpdatedAt: time.Now(),
// 	// 		// 		})
// 	// 		// 	}
// 	// 			// user, _ := GetFacebookProfile(fbID)

// 	// 			// if err != nil {

// 	// 			// 	// if err != nil {
// 	// 			// 	// 	log.Println("get user error", err)
// 	// 			// 	// 	continue
// 	// 			// 	// }
// 	// 			// 	gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).UpdateColumns(&models.Customer{
// 	// 			// 		Image:     user.ProfilePic,
// 	// 			// 		FirstName: user.FirstName,
// 	// 			// 		LastName:  user.LastName,
// 	// 			// 	})

// 	// 			// }

// 	// 			// if senderID == "" || recipientID == "" {
// 	// 			// 	log.Println("Skipping message with empty sender or recipient ID")
// 	// 			// 	continue
// 	// 			// }

// 	// 			// // err := dbservice.CreateMesseng(&models.Chat{
// 	// 			// // 	SenderID:    senderID,
// 	// 			// // 	UserID:      "1e55b100-8a4e-4372-a9e9-7d3c5f4a2a77", // You might want to dynamically look up user ID
// 	// 			// // 	RecipientID: recipientID,
// 	// 			// // 	JSONMesseng: string(c.BodyRaw()),
// 	// 			// // })

// 	// 			// gormpkg.GetDB().Table(models.TableNameCustomer).Where("facebook_id = ?", fbID).First(&user)

// 	// 			// var payload interface{}
// 	// 			// json.Unmarshal(c.BodyRaw(), &payload)

// 	// 			//  UpdateColumns(&models.Customer{
// 	// 			// 		FirstName: fullnam,
// 	// 			// 	})
// 	// 			// PushToUser(fbID, fiber.Map{
// 	// 			// 	// "customer_id": fbID,
// 	// 			// 	"user":    user,
// 	// 			// 	"message": payload,
// 	// 			// })

// 	// 			// PushToAll(fiber.Map{
// 	// 			// 	// "customer_id": fbID,
// 	// 			// 	"user":    user,
// 	// 			// 	"message": payload,
// 	// 			// })

// 	// 			// if err != nil {
// 	// 			// 	log.Println("Failed to create message:", err)
// 	// 			// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 	// 			// 		"error": "create messeng",
// 	// 			// 	})
// 	// 			// }

// 	// 			// Handle postbacks
// 	// 			// if msg.Postback != nil && msg.Postback.Payload != "" {
// 	// 			// 	log.Printf("Received postback from %s: %s\n", senderID, msg.Postback.Payload)
// 	// 			// 	SendMessage(senderID, "You clicked: "+msg.Postback.Payload)
// 	// 			// }

// 	// 			// // Handle delivery confirmations
// 	// 			// if msg.Delivery != nil {
// 	// 			// 	log.Printf("Delivery confirmed for %d message(s): %v\n",
// 	// 			// 		len(msg.Delivery.Mids), msg.Delivery.Mids)
// 	// 			// }
// 	// 		}
// 	// 	}
// 	// }

// 	return c.SendStatus(fiber.StatusOK)
// }

// // --- Configuration ---
// // This secret key must match the one configured on the sending service's side.
// const webhookSecret = "EAARDcwZBMbeQBPZCG1ZAHM1x"

// // Define a struct to match the expected JSON payload
// type MessageUpdate struct {
// 	EventType string `json:"event_type"`
// 	MessageID string `json:"message_id"`
// 	Status    string `json:"status"` // e.g., "delivered", "read", "failed"
// }

// --- Handler Function ---

func FacebookWebhookHandler(c *fiber.Ctx) error {

	// --- 1. Handle Verification Request (GET) ---
	if c.Method() == fiber.MethodGet {
		mode := c.Query("hub.mode")
		// token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		if mode == "subscribe" {
			log.Println("✅ Webhook verified successfully by Facebook.")
			// Respond with the hub.challenge value
			return c.Status(fiber.StatusOK).SendString(challenge)
		}
		return c.Status(fiber.StatusForbidden).SendString("Forbidden: Invalid verification parameters")
	}

	// --- 2. Handle Event Notification (POST) ---
	if c.Method() == fiber.MethodPost {
		// a. Get raw body for signature validation
		// payloadBody := c.Body()

		// b. Validate the Signature
		// Facebook uses X-Hub-Signature-256 header
		providedSignature := c.Get("X-Hub-Signature-256")
		if providedSignature == "" {
			log.Println("Security validation failed: Missing X-Hub-Signature-256 header")
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Missing signature")
		}

		// if !validateFacebookSignature(payloadBody, providedSignature, webhookSecret) {
		// 	log.Println("Security validation failed: Signature mismatch")
		// 	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid signature")
		// }
		// log.Println("Signature validated successfully.")

		// c. Decode and Process the Payload
		// The Facebook payload is complex and often contains 'entry' arrays.

		// var webhookData map[string]interface{}
		// if err := c.BodyParser(&webhookData); err != nil {
		// 	log.Printf("Error decoding JSON payload: %v", err)
		// 	return c.Status(fiber.StatusBadRequest).SendString("Bad Request: Invalid JSON format")
		// }

		// // -----------------------------------------------------------------
		// //  วิธีพิมพ์ webhookData ทั้งหมด (แบบอ่านง่าย)
		// // -----------------------------------------------------------------
		// jsonData, err := json.MarshalIndent(webhookData, "", "  ") // "  " คือการย่อหน้า
		// if err != nil {
		// 	log.Printf("Error marshaling webhook data to JSON: %v", err)
		// } else {
		// 	// พิมพ์ข้อมูล JSON ทั้งหมดที่ได้รับออกมา
		// 	log.Printf("--- Received Full Webhook Payload --- :\n%s", string(jsonData))
		// }

		// c. Decode Payload
		// เปลี่ยนจาก map[string]interface{} เป็น Struct FacebookPayload
		var payload FacebookPayload
		if err := c.BodyParser(&payload); err != nil {
			log.Printf("Error decoding JSON payload: %v", err)
			return c.Status(fiber.StatusBadRequest).SendString("Bad Request: Invalid JSON format")
		}

		// (Optional) Print JSON แบบอ่านง่ายเหมือนเดิม
		jsonData, _ := json.MarshalIndent(payload, "", "  ")
		log.Printf("--- Received Webhook --- :\n%s", string(jsonData))

		// -----------------------------------------------------------------
		//  Logic: Check Ad ID and Save to DB
		// -----------------------------------------------------------------

		// 1. Loop เข้าไปใน entry
		for _, entry := range payload.Entry {
			// 2. Loop เข้าไปใน messaging
			for _, event := range entry.Messaging {

				// ตรวจสอบว่าเป็น Message และมี Referral (มาจากโฆษณา) หรือไม่
				if event.Message != nil && event.Message.Referral != nil {

					// ดึงค่า Ad ID
					adID := event.Message.Referral.AdID

					// ตรวจสอบว่ามี Ad ID จริงๆ ไม่ใช่ค่าว่าง
					if adID != "" {
						userID := event.Sender.ID
						// text := event.Message.Text
						// [เพิ่ม]: ดึง Ad Title (ต้องเช็ค nil ก่อนเผื่อไม่มี)
						adTitle := ""
						if event.Message.Referral.AdsContextData != nil {
							adTitle = event.Message.Referral.AdsContextData.AdTitle
						}

						log.Printf("🎯 Found Ad Click! Title: %s (ID: %s)", adTitle, adID)
						log.Printf("🎯 Found Ad Click! User: %s, Ad ID: %s", userID, adID)

						// // บันทึกลง Database
						// saveAdTrackingToDB(userID, adID, text)
					}
				}
			}
		}

		// -----------------------------------------------------------------

		// Example logging the object type (page, user, etc.)
		// log.Printf("Received Facebook Webhook Event: Object Type: %v", webhookData["object"])

		// TODO: Implement parsing and processing of the nested 'entry' and 'messaging' arrays here.

		// d. Send Success Response
		// Must return 200 OK quickly, even before processing is complete.
		return c.SendStatus(fiber.StatusOK)
	}

	// Handle other methods
	return c.SendStatus(fiber.StatusMethodNotAllowed)
}

// --- Security Helper ---

// validateFacebookSignature verifies the HMAC signature using SHA256.
// func validateFacebookSignature(payloadBody []byte, providedSignature string, secret string) bool {
// 	// Facebook signature is prefixed with 'sha256='
// 	const signaturePrefix = "sha256="
// 	if len(providedSignature) < len(signaturePrefix) || providedSignature[:len(signaturePrefix)] != signaturePrefix {
// 		return false // Prefix not found
// 	}

// 	// Strip the prefix
// 	providedSignature = providedSignature[len(signaturePrefix):]

// 	key := []byte(secret)
// 	h := hmac.New(sha256.New, key)
// 	h.Write(payloadBody)
// 	expectedMAC := hex.EncodeToString(h.Sum(nil))

// 	// Use constant-time comparison for security
// 	return hmac.Equal([]byte(expectedMAC), []byte(providedSignature))
// }
