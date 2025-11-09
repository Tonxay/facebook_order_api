package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

// !! ใส่ App ID และ Secret ของคุณตรงนี้ !!
const (
	FACEBOOK_APP_ID     = "YOUR_APP_ID"
	FACEBOOK_APP_SECRET = "YOUR_APP_SECRET"
)

// --- Structs (เหมือนเดิม) ---
type AuthRequest struct {
	AccessToken string `json:"accessToken"`
}
type LongLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type Page struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}
type PagesResponse struct {
	Data []Page `json:"data"`
}

// --- Handler (เวอร์ชัน Fiber) ---
// สังเกตว่าเปลี่ยนจาก (w, r) เป็น (c *fiber.Ctx)
func FacebookCallbackHandler(c *fiber.Ctx) error {

	// // 2. อ่าน JSON (Fiber ทำได้ง่ายกว่า)
	// var req AuthRequest
	// if err := c.BodyParser(&req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"status": "error", "message": "Invalid request body",
	// 	})
	// }
	// shortLivedToken := req.AccessToken

	// // 3. แลก Token (เรียกฟังก์ชันเดิม)
	// longLivedToken, err := exchangeToken(shortLivedToken)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status": "error", "message": err.Error(),
	// 	})
	// }
	// log.Printf("ได้ Long-Lived User Token มาแล้ว: %s...", longLivedToken[:10])

	// // 4. ดึง Page Token (เรียกฟังก์ชันเดิม)
	// pages, err := getPages(longLivedToken)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status": "error", "message": err.Error(),
	// 	})
	// }

	// // 5. บันทึกผลลัพธ์ (เหมือนเดิม)
	// log.Println("--- พบเพจที่เชื่อมต่อ ---")
	// for _, page := range pages.Data {
	// 	log.Printf("  Page ID: %s", page.ID)
	// 	log.Printf("  Page Name: %s", page.Name)
	// 	log.Printf("  Page Access Token: %s...", page.AccessToken[:10])
	// 	// **** จุดนี้คือจุดที่คุณต้องบันทึก page.AccessToken ลง Database ****
	// }
	// log.Println("--------------------------")

	// 6. ตอบกลับ (Fiber ทำได้ง่ายกว่า)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("เชื่อมต่อสำเร็จ เพจ!"),
	})
}

// --- Helper Functions (เหมือนเดิมเป๊ะ) ---
// (ฟังก์ชัน 2 อันนี้ไม่ต้องแก้เลย)

func exchangeToken(shortLivedToken string) (string, error) {
	apiURL := "https://graph.facebook.com/v19.0/oauth/access_token"
	params := url.Values{}
	params.Set("grant_type", "fb_exchange_token")
	params.Set("client_id", FACEBOOK_APP_ID)
	params.Set("client_secret", FACEBOOK_APP_SECRET)
	params.Set("fb_exchange_token", shortLivedToken)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResponse LongLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("error decoding token response: %s", string(body))
	}

	if tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("did not get access token: %s", string(body))
	}

	return tokenResponse.AccessToken, nil
}

func getPages(longLivedUserToken string) (*PagesResponse, error) {
	apiURL := "https://graph.facebook.com/v19.0/me/accounts?fields=id,name,access_token&access_token=" + longLivedUserToken

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var pagesResponse PagesResponse
	if err := json.Unmarshal(body, &pagesResponse); err != nil {
		return nil, fmt.Errorf("error decoding pages response: %s", string(body))
	}

	return &pagesResponse, nil
}

// // --- Server (เวอร์ชัน Fiber) ---
// func main() {
// 	// 1. สร้างแอป Fiber
// 	app := fiber.New()

// 	// 2. ใช้ Middleware (CORS) ง่ายๆ แค่บรรทัดเดียว
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "*", // (ใน Production ควรระบุโดเมนจริง)
// 		AllowMethods: "POST, OPTIONS",
// 	}))

// 	// 3. กำหนด Route (เส้นทาง)
// 	app.Post("/auth/facebook/callback", facebookCallbackHandler)

// 	// 4. สตาร์ทเซิร์ฟเวอร์
// 	log.Println("Go Fiber Backend Server กำลังเริ่มที่ http://localhost:8080")
// 	if err := app.Listen(":8080"); err != nil {
// 		log.Fatal(err)
// 	}
// }
