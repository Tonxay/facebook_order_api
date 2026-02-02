package service

import (
	"github.com/go-resty/resty/v2"
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
