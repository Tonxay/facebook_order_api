package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gofiber/fiber/v2"
)

// โครงสร้างข้อมูลสำหรับรับจาก TikWM API
type TikWMResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Title string `json:"title"`
		Play  string `json:"play"` // ลิงก์ No Watermark
		Cover string `json:"cover"`
		Music string `json:"music"`
	} `json:"data"`
}

// TictokVideoService คือ Handler สำหรับจัดการการดาวน์โหลด
func TictokVideoService(c *fiber.Ctx) error {
	// 1. รับค่า URL จาก Query Parameter (?url=...)
	tiktokURL := c.Query("url")

	// 2. Validation เบื้องต้น
	if tiktokURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "กรุณาส่งพารามิเตอร์ url มาด้วย",
		})
	}

	// 3. เรียกใช้ API ภายนอก (TikWM)
	client := resty.New()
	var tikRes TikWMResponse

	_, err := client.R().
		SetQueryParam("url", tiktokURL).
		SetResult(&tikRes).
		Get("https://www.tikwm.com/api/")

	// 4. ตรวจสอบ Error จากการ Request หรือจาก API
	if err != nil || tikRes.Code != 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลวิดีโอได้ โปรดตรวจสอบ URL อีกครั้ง",
		})
	}

	// 5. ส่งผลลัพธ์กลับไปให้ Client
	// หมายเหตุ: TikWM จะส่ง path ของวิดีโอมา เราต้องเติม domain ข้างหน้า
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"title":       tikRes.Data.Title,
			"video_no_wm": tikRes.Data.Play,
			"cover":       tikRes.Data.Cover,
			"music":       tikRes.Data.Music,
		},
	})
}

func PinduoduoService(c *fiber.Ctx) error {
	targetURL := c.Query("url")
	if targetURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "URL is required"})
	}

	// 1. ตั้งค่า Browser (ปลอมตัวเป็น Mobile)
	l := launcher.New().
		Headless(true). // เปิดแบบไม่โชว์หน้าต่าง
		Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/04.1")

	browserURL := l.MustLaunch()
	browser := rod.New().ControlURL(browserURL).MustConnect()
	defer browser.MustClose()

	// 2. เปิดหน้าเว็บและรอให้วิดีโอโหลด
	page := browser.MustPage(targetURL)

	// รอจนกว่าแท็ก video จะปรากฏ (timeout 10 วินาที)
	var videoSrc string
	err := rod.Try(func() {
		page.MustWaitLoad()
		// ค้นหาแท็ก video และดึง src
		el := page.MustElement("video")
		videoSrc = ""
		if src := el.MustAttribute("src"); src != nil {
			videoSrc = *src
		}
		if videoSrc == "" {
			src := el.MustElement("source").MustAttribute("src")
			if src != nil {
				videoSrc = *src
			}
		}
	})

	if err != nil || videoSrc == "" {
		return c.Status(500).JSON(fiber.Map{"error": "หาลิงก์วิดีโอไม่เจอ (อาจต้องใช้การวิเคราะห์ขั้นสูง)"})
	}

	// 3. ส่งข้อมูลกลับไปให้ Flutter
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"title":       "Pinduoduo Video",
			"video_no_wm": videoSrc,
			"cover":       "", // สามารถใช้ rod ดึงรูป og:image เพิ่มได้
		},
	})
}
