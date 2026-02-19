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
	// token := AnousithLoging("92339355", "s0987654")
	// token := AnousithLoging("2092339355", "s0987654")
	order, _ := dbservice.GetOrders(gormpkg.GetDB(), request.StatusOrderRequest{
		IsCancel:   false,
		ShippingID: "3531696c-af5b-4d73-b677-dbd9e1bb4f1b",
		Statuses:   []string{"shipped"},
	})

	if len(order) == 0 {
		return fmt.Errorf("No orders found")
	}

	amountBills := HardCodeAnousithOrder()

	//AnousithOrder(token)

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
	// 1. ป้องกันค่าว่าง (สาเหตุที่ทำให้โดนแบน ResultLimitToDay)
	cleanTel := strings.TrimSpace(tel)
	cleanPass := strings.TrimSpace(password)
	if cleanTel == "" || cleanPass == "" {
		log.Println("❌ ยกเลิกการยิง API: เบอร์โทรหรือรหัสผ่านเป็นค่าว่าง")
		return ""
	}

	payload := map[string]interface{}{
		"operationName": "CustomerLogin",
		"variables": map[string]interface{}{
			"where": map[string]interface{}{
				"username": cleanTel,
				"password": cleanPass,
			},
		},
		"query": "mutation CustomerLogin($where: CustomerLoginInput!) {\n  customerLogin(where: $where) {\n    accessToken\n    data {\n      id_list\n      full_name\n      profile_img\n      status\n      contact_info\n      address\n      village\n      district {\n        id_list\n        title\n      }\n      state {\n        provinceName\n        id_state\n      }\n      Bank_KIP\n      BANK_THB\n      BANK_USD\n      BANK_NAME\n      gender\n      isActive\n      isVerify\n    }\n  }\n}",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error creating JSON payload:", err)
		return "" // เปลี่ยนจาก log.Fatal เพื่อไม่ให้ POS ล่ม
	}

	apiUrl := "https://pro.api.anousith.express/graphql"

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating Request:", err)
		return ""
	}

	// 2. ตั้งค่า Headers และบังคับล้าง Cache (วางถูกตำแหน่งแล้ว!)
	req.Close = true // บังคับให้ตัด Connection หลังยิงเสร็จ (ลบความจำเซิร์ฟเวอร์)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("referer", "https://app.anousith.express/")
	req.Header.Set("Cache-Control", "no-cache") // สั่งเซิร์ฟเวอร์ว่าห้ามเอาของเก่ามาตอบ

	// 3. สร้าง Client ส่วนตัวขึ้นมาใหม่เลย (เลิกใช้ DefaultClient เด็ดขาด)
	myClient := &http.Client{
		Timeout: 15 * time.Second, // ป้องกันแอปค้างถ้ายิงไม่เข้า
	}

	// ยิง Request ด้วย myClient
	resp, err := myClient.Do(req)
	if err != nil {
		log.Println("HTTP Request Failed:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return ""
	}

	// ตรวจสอบ Error จากเซิร์ฟเวอร์
	var respJSON map[string]interface{}
	if err := json.Unmarshal(body, &respJSON); err != nil {
		log.Println("Error parsing response JSON:", err)
		return ""
	}

	if errors, ok := respJSON["errors"].([]interface{}); ok && len(errors) > 0 {
		log.Printf("⚠️ เซิร์ฟเวอร์ตอบกลับเป็น Error: %v\n", errors[0])
		return ""
	}

	// ดึง Token
	accessToken := ""
	if data, ok := respJSON["data"].(map[string]interface{}); ok {
		if cl, ok := data["customerLogin"].(map[string]interface{}); ok {
			if at, ok := cl["accessToken"].(string); ok {
				accessToken = at
			}
		}
	}

	if accessToken == "" {
		log.Printf("ไม่พบ accessToken ใน Response: %s", string(body))
	} else {
		log.Println("✅ ล็อกอินสำเร็จ!")
	}

	return accessToken
}

func HardCodeAnousithOrder() custommodel.ItemsV2Response {
	rawJSON := `{
    "data": {
        "itemsV2": {
            "total": 67,
            "data": [
                {
                    "_id": "296283817",
                    "trackingPlatform": "8260503410231",
                    "trackingId": "8260503410231",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ຈັນທາແສນຄຳ",
                    "receiverPhone": "56556282",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:50:20.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:24:31.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ໄຊຍະບູລີ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເມືອງຊຽງຮ່ອນ",
                        "branch_address": "ບ້ານ ໃໝ່ໜອງຊາງ",
                        "contactInfo": "59751244 ຢູ່ຕິດກັບ ໂຮງຮຽນ ມສ ຊຽງຮ່ອນ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283769",
                    "trackingPlatform": "8260501012640",
                    "trackingId": "8260501012640",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ຈັນທະຈອນ ລັດດາ",
                    "receiverPhone": "91845464",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:34.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:23:53.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຄຳມ່ວນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາໄຊບົວທອງ(ເມືອງໄຊບົວທອງ)",
                        "branch_address": "ບ້ານ ສີວິໄລ",
                        "contactInfo": "58438885 95855512"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283743",
                    "trackingPlatform": "8260502602788",
                    "trackingId": "8260502602788",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Kongla",
                    "receiverPhone": "58095421",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:59.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:23:29.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງນ້ຳທາ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ 4ແຍກໂພນໄຊ(ຫຼວງນໍ້າທາ)",
                        "branch_address": "ບ້ານ ໂພນໄຊ",
                        "contactInfo": "02051611766"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283720",
                    "trackingPlatform": "8260504489498",
                    "trackingId": "8260504489498",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Phasaphone",
                    "receiverPhone": "52144045",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:29.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:23:06.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຈອມທອງ(ເມືອງວຽງທອງ)",
                        "branch_address": "ບ້ານ ຈອມທອງ",
                        "contactInfo": "29925353"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283673",
                    "trackingPlatform": "8260502859853",
                    "trackingId": "8260502859853",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ສີເມື່ອງ",
                    "receiverPhone": "99477155",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:24.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:22:25.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານກາງຫຼັກ5",
                        "branch_address": "ບ້ານ ກາງຫຼັກ5",
                        "contactInfo": " ❤❤    99609586  96424926   78801844 ❤❤"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283604",
                    "trackingPlatform": "8260508114514",
                    "trackingId": "8260508114514",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ພີອອດ",
                    "receiverPhone": "96405512",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:10.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:21:43.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ປາກແບງ(ເມືອງ ປາກແບງ)",
                        "branch_address": "ບ້ານ ປາກແບງ",
                        "contactInfo": "ເບີຕິດຕໍ່ ປະຈຳສາຂາ 91823221"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283598",
                    "trackingPlatform": "8260508409347",
                    "trackingId": "8260508409347",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ເສືອພັນທະສູນທອນ",
                    "receiverPhone": "99339869",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:13.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:21:37.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຂີ້ນາກ2",
                        "branch_address": "ບ້ານ ຂີ້ນາກ",
                        "contactInfo": "95395194/78320312"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283545",
                    "trackingPlatform": "8260507693647",
                    "trackingId": "8260507693647",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Modlnthongson",
                    "receiverPhone": "78092692",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:51.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:20:55.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫຼັງຕະຫຼາດໃໝ່(ປາກຊັນ)",
                        "branch_address": "ບ້ານ ໂພນໄຊ",
                        "contactInfo": "0309924266/99690841/57369597 ຕິດຕໍ່ພະນັກງານ,(95323661  ເຈົ້າຂອງສາຂາຕິດສະເພາະມີບັນຫາກ່ຽວກັບພະນັກງານ)"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283513",
                    "trackingPlatform": "8260503615507",
                    "trackingId": "8260503615507",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ກິນ້ອຍ",
                    "receiverPhone": "98097731",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:21.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:20:29.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໜອງກອກ ( ບາຈຽງ )",
                        "branch_address": "ບ້ານ ໜອງກອກ",
                        "contactInfo": "76273334"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283472",
                    "trackingPlatform": "8260503485274",
                    "trackingId": "8260503485274",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ອະນຸສອນ",
                    "receiverPhone": "96771922",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:15.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:19:58.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານວັງເຕົ່າ",
                        "branch_address": "ບ້ານ ວັງເຕົ່າ",
                        "contactInfo": "93333780 ເບີໂທແລະວັອດແອັບ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283458",
                    "trackingPlatform": "8260507030063",
                    "trackingId": "8260507030063",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ຈັນທະໝອນ",
                    "receiverPhone": "96912345",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:03.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:19:50.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຊຽງຂວາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ລາດງ່ອນ",
                        "branch_address": "ບ້ານລາດງ່ອນ",
                        "contactInfo": "54148009/52588277"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283412",
                    "trackingPlatform": "8260504595537",
                    "trackingId": "8260504595537",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Monthidakone",
                    "receiverPhone": "98813846",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:06.000Z",
                    "width": 58,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:19:17.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສະຫວັນນະເຂດ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫຼຽນໄຊ(ອາດສະພັງທອງ)",
                        "branch_address": "ບ້ານ ຫຼຽນໄຊ",
                        "contactInfo": "ັ9267668/ສາຍດ່ວນ92676668/wpp95448143"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283392",
                    "trackingPlatform": "8260500012531",
                    "trackingId": "8260500012531",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ລັກຊະນະ",
                    "receiverPhone": "99959539",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:08.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:19:03.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຮ່ອງຂະຍອມ",
                        "branch_address": "ບ້ານ ຮ່ອງຂະຍອມ",
                        "contactInfo": " ກົງກັນຂ້າມກັບຮ້ານເຂົ້າປຸ້ນນ້ຳແຈ່ວ, ເສັ້ນທາງໄປໂບ້ລິ້ງ, ໂທ0305472258 /ແອ໋ບ94472110/ໂທ 02091734348"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283346",
                    "trackingPlatform": "8260507348444",
                    "trackingId": "8260507348444",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Phong Phong phun",
                    "receiverPhone": "99593249",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:41.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:18:23.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ໄຊຍະບູລີ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ນໍ້າຊົ້ງ",
                        "branch_address": "ບ້ານ ນໍ້າຊົ້ງ",
                        "contactInfo": "99425558"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283334",
                    "trackingPlatform": "8260509188155",
                    "trackingId": "8260509188155",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Bobby",
                    "receiverPhone": "52520505",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:33.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:18:14.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານເມືອງງາ(ນະຄອນຫຼວງພະບາງ)",
                        "branch_address": "ບ້ານ ເມືອງງາ",
                        "contactInfo": "(ແຈ້ງເລື່ອງທົ່ວໄປ28622444),(ຍ້າຍບີນຕິດຕໍ່ເບີ 99915592)"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283272",
                    "trackingPlatform": "8260509648915",
                    "trackingId": "8260509648915",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Kingphet",
                    "receiverPhone": "58990809",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:17.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:17:28.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໂນນສະຫວັນ(ເມືອງໂພນທອງ)",
                        "branch_address": "ບ້ານ ໂນນສະຫວັນ",
                        "contactInfo": "99977694/92266448ເບີປະຈຳຮ້ານ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283267",
                    "trackingPlatform": "8260503064158",
                    "trackingId": "8260503064158",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Ketsadaphone",
                    "receiverPhone": "54888230",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:43.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:17:25.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສະຫວັນນະເຂດ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ວົງວຽນເຊໂນ",
                        "branch_address": "ບ້ານ ເຊໂນ",
                        "contactInfo": "/56677682/ /92983181/  98667458 / 8:30ຫາ17:00"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283176",
                    "trackingPlatform": "8260507773199",
                    "trackingId": "8260507773199",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "see lea 99517049",
                    "receiverPhone": "99517049",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:33.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:16:20.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສາລະວັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຄົງນະຄອນ(ຄົງເຊໂດນ)",
                        "branch_address": "ບ້ານ ຄົງນະຄອນ",
                        "contactInfo": "99423292"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283093",
                    "trackingPlatform": "8260508731716",
                    "trackingId": "8260508731716",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "xayavong",
                    "receiverPhone": "98249800",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:37.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:15:25.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໂຮມສຸກ(ເມືອງໄຊ)",
                        "branch_address": "ບ້ານ ໂຮມສຸກ",
                        "contactInfo": "28722425/0309440559"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296283086",
                    "trackingPlatform": "8260507267345",
                    "trackingId": "8260507267345",
                    "itemName": "ຄ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Noysihavong",
                    "receiverPhone": "96311799",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:32.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:15:21.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາດອນແກ້ວ(ອຸດົມໄຊ)",
                        "branch_address": "ບ້ານ ດອນແກ້ວ",
                        "contactInfo": "WhataApp ພະນັກງານປະຈຳສາຂາດອນແກ້ວ  02099358172"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282989",
                    "trackingPlatform": "8260501377432",
                    "trackingId": "8260501377432",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ການຟອຍ",
                    "receiverPhone": "76813288",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:29.000Z",
                    "width": 58,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:14:19.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫຼັກ20 (ບ້ານ ຫ້ວຍແກ້ວ ເມືອງ ຄຳເກີດ)",
                        "branch_address": "ບ້ານ ຫວ້ຍແກ້ວ",
                        "contactInfo": "ໂທ: 0305637163,ວ໊ອດແອັບ: 02097636162"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282869",
                    "trackingPlatform": "8260504043611",
                    "trackingId": "8260504043611",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 150000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ບຸນແກ້ວ",
                    "receiverPhone": "59527878",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:17.000Z",
                    "width": 58,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:13:04.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫົວພັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເມືອງຮາມ(ບ້ານ ຮາມໃຕ້) ເມືອງ ຊຳເໜືອ",
                        "branch_address": "ບ້ານ ຮາມໃຕ້",
                        "contactInfo": "02055664615/02058308333"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282765",
                    "trackingPlatform": "8260507110065",
                    "trackingId": "8260507110065",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ສາລິລາຍ",
                    "receiverPhone": "99737463",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:38.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:11:57.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອັດຕະປື"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຝັ່ງແດງ(ເມືອງ ໄຊເສດຖາ)",
                        "branch_address": "ບ້ານ ຝັ່ງແດງ",
                        "contactInfo": " ເບີຕິດຕໍ່ວຽກຕ່າງໆເບີແອັບກັບໂທ93554575"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282745",
                    "trackingPlatform": "8260503775256",
                    "trackingId": "8260503775256",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Lar phalaphon",
                    "receiverPhone": "97969565",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:29.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:11:44.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໂພນທອງ",
                        "branch_address": "ບ້ານ ໂພນທອງ",
                        "contactInfo": "91438384ວ໊ອດແອັບ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282669",
                    "trackingPlatform": "8260503592055",
                    "trackingId": "8260503592055",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Gick",
                    "receiverPhone": "95023820",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:24.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:11:03.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ດົງປ່າລານ",
                        "branch_address": "ບ້ານດົງປ່າລານ",
                        "contactInfo": "22659159/91199112ເປີດ8:30ປິດ17:30ຮ່ອມ ສໍລະດິດ ເລື້ອງກົດCODແມ່ນສົງຫາສາຂາກ່ອນ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282646",
                    "trackingPlatform": "8260507065850",
                    "trackingId": "8260507065850",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Jimmy chanthajon",
                    "receiverPhone": "28529240",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:41.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:10:50.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສາລະວັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ນາເຫຼັກ(ສາລະວັນ)",
                        "branch_address": "ບ້ານ ນາເຫຼັກ",
                        "contactInfo": "02096795878  .0305472181 "
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282593",
                    "trackingPlatform": "8260506290348",
                    "trackingId": "8260506290348",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Pataiy",
                    "receiverPhone": "95669283",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:26.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:10:19.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ທົ່ງກາງ (ເຂດທ່າພະລານໄຊ)",
                        "branch_address": "ບ້ານ ທ່າພະລານໄຊ",
                        "contactInfo": "(98019732ເບີສາຂາ),(94056951 ຍ້າຍ ແກ້ໄຂ າ)ເຂດທ່າພະລານໄຊ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282538",
                    "trackingPlatform": "8260506143912",
                    "trackingId": "8260506143912",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ນິກຄົນເດີມ",
                    "receiverPhone": "91486991",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:46.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:09:49.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ໄຊຍະບູລີ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານຂອນ(ເມືອງເງິນ)",
                        "branch_address": "ບ້ານ ຂອນ",
                        "contactInfo": "54869948"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282469",
                    "trackingPlatform": "8260507397577",
                    "trackingId": "8260507397577",
                    "itemName": "ຄຂ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Samlansamoan",
                    "receiverPhone": "55424741",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:21.000Z",
                    "width": 57,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:09:07.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "branch_address": "ບ້ານ ວຽງຈະເລີນ",
                        "contactInfo": "02055966379"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282446",
                    "trackingPlatform": "8260506737513",
                    "trackingId": "8260506737513",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "loun lpb",
                    "receiverPhone": "77958424",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:46.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:08:56.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ປາກແວດ(ເມືອງ ຊຽງເງິນ)",
                        "branch_address": "ບ້ານ ປາກແວດ",
                        "contactInfo": "whatsapp:79996611, ໂທ 92008033 , ພິກັດ: ຕໍ່ໜ້າຕະຫຼາດແລງບ້ານປາກແວດ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282341",
                    "trackingPlatform": "8260504540504",
                    "trackingId": "8260504540504",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ພອນໄຊ",
                    "receiverPhone": "92361995",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:56.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:07:53.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສະຫວັນນະເຂດ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເຊທ່າມວກ",
                        "branch_address": "ບ້ານ ເຊທ່າມວກ",
                        "contactInfo": "020 94476565/020 99756665 ຖ້າສົ່ງຂໍ້ຄວາມບໍ່ຕອບໂທເອົາເດີ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282318",
                    "trackingPlatform": "8260501448999",
                    "trackingId": "8260501448999",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ຊາຍວັນນະສີນ",
                    "receiverPhone": "93377233",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:15.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:07:39.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສາລະວັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເລົ່າງາມ",
                        "branch_address": "ບ້ານ ເລົ່າງາມ",
                        "contactInfo": "ແຈ້ງຍ້າຍບິນ ແຈ້ງເກັບCOD ຕິດຕໍ່ເບີ W.A 02023434451"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282234",
                    "trackingPlatform": "8260501378995",
                    "trackingId": "8260501378995",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Khamfong",
                    "receiverPhone": "55434598",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:28.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:06:57.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫຼັກ4 ( ປາກເຊ)",
                        "branch_address": "ບ້ານ ພູມ່ວງ",
                        "contactInfo": "97452420 ມີບໍລິການຈັດສົ່ງເຄື່ອງເຖິງໜ້າບ້ານ, ສາມາດສອບຖາມເງື່ອນໄຂການຂົນສົ່ງເພີ່ມເຕີມໄດ້."
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282226",
                    "trackingPlatform": "8260505173151",
                    "trackingId": "8260505173151",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Xaiphone",
                    "receiverPhone": "56940131",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:40:38.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:06:54.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານດອນ(ເມືອງຫຼາ)",
                        "branch_address": "ບ້ານ ດອນ",
                        "contactInfo": "92847696"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282151",
                    "trackingPlatform": "8260504286325",
                    "trackingId": "8260504286325",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Paensvy",
                    "receiverPhone": "91482814",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:21.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:06:13.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານຫ້ວຍເລີກ(ເມືອງທ່າພະບາດ)",
                        "branch_address": "ບ້ານ ຫ້ວຍເລີກ",
                        "contactInfo": "29938606"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296282063",
                    "trackingPlatform": "8260501013637",
                    "trackingId": "8260501013637",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Paoyanong",
                    "receiverPhone": "91948471",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:18.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:05:16.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງນ້ຳທາ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບໍເຕັນ(ເມືອງ ຫຼວງນໍ້າທາ)",
                        "branch_address": "ບ້ານ ບໍເຕັນ",
                        "contactInfo": "52017183/78830607/99554358/97304611ບ້ານບໍ່ເຕັນ "
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281989",
                    "trackingPlatform": "8260507713343",
                    "trackingId": "8260507713343",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ສົມຫວັງ",
                    "receiverPhone": "78492219",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:10.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:04:35.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໜອງບົວທອງເໜືອ(ເຂດໜອງປິງ)",
                        "branch_address": "ບ້ານ ໜອງບົວທອງເໜືອ",
                        "contactInfo": "ປ່ຽນເບີໃຫມ່ 59979546 ຢູ່ແຖວ3ແຍກທາງລົງໜອງປີງ,ຕິດກັບຊຸບເປີມາເກັດໂຊກທະວີ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281922",
                    "trackingPlatform": "8260503973937",
                    "trackingId": "8260503973937",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "sipaser",
                    "receiverPhone": "56228721",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:06.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:03:58.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ທ່າຫີນ ( ສະພານໄຊ )",
                        "branch_address": "ບ້ານ ສະພານໄຊ",
                        "contactInfo": "ສາຂາ ທ່າຫີນ 91207692 ພະນັກງານ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281807",
                    "trackingPlatform": "8260504796283",
                    "trackingId": "8260504796283",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ໜອດລັດສະໝີ",
                    "receiverPhone": "98952425",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:14.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:02:55.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານວັງ(ເມືອງ ໝື່ນ)",
                        "branch_address": "ບ້ານ ວັງ",
                        "contactInfo": "020 59391473"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281665",
                    "trackingPlatform": "8260502151092",
                    "trackingId": "8260502151092",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "fangko",
                    "receiverPhone": "96070078",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:36.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T14:01:37.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານຫຼັກ36(ເມືອງປາກຊ່ອງ)",
                        "branch_address": "ບ້ານ ຫຼັກ36",
                        "contactInfo": "58868325  59982902"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281291",
                    "trackingPlatform": "8260505364728",
                    "trackingId": "8260505364728",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "sab bkt",
                    "receiverPhone": "99101700",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:10.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:58:07.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫົວພັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຊອນເໜືອ",
                        "branch_address": "ບ້ານ ຊອນເໜືອ",
                        "contactInfo": "020 92477760; 020 92265586"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281212",
                    "trackingPlatform": "8260504243539",
                    "trackingId": "8260504243539",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Anousit pongvilay",
                    "receiverPhone": "58134354",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:42:56.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:57:24.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສາລະວັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ແບ່ງດ່ານ(ເມືອງລະຄອນເພັງ)",
                        "branch_address": "ບ້ານ ແບ່ງດ່ານ",
                        "contactInfo": "95568445ເຈົ້າຂອງສາຂາ,98179523ພງປະຈຳສາສາຂາ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296281076",
                    "trackingPlatform": "8260507662678",
                    "trackingId": "8260507662678",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 160000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Ler hackler",
                    "receiverPhone": "92555508",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:23.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:56:18.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ດອນງີ້ວ(ເມືອງງາ)",
                        "branch_address": "ບ້ານ ດອນງີ້ວ",
                        "contactInfo": "ແກ້ໄຂບິນ02052066440/ຍ້າຍສາຂາ02099578263"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280936",
                    "trackingPlatform": "8260503142288",
                    "trackingId": "8260503142288",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ວົງສະວັນ ອີນມະນີ",
                    "receiverPhone": "55196677",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:04.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:55:08.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເມືອງວຽງຄຳ(ຫຼວງພະບາງ)",
                        "branch_address": "ບ້ານ ວຽງຄຳ",
                        "contactInfo": " 91792167/ 97998818"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280797",
                    "trackingPlatform": "8260503317358",
                    "trackingId": "8260503317358",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Vilay xayyaveth",
                    "receiverPhone": "55654498",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:48.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:54:02.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ວັດນາກ(ສີສັດຕະນາກ)",
                        "branch_address": "ບ້ານ ວັດນາກ",
                        "contactInfo": " ( ໃກ້ກັບ ວິທະຍາໄລ ການຊ່າງລາວ-ເຢຍລະມັນ) whatsapp  55800552"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280635",
                    "trackingPlatform": "8260509929870",
                    "trackingId": "8260509929870",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ຈິດຕະພອນ",
                    "receiverPhone": "91001659",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:32.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:52:47.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງນ້ຳທາ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານວຽງເໜືອ",
                        "branch_address": "ບ້ານ ວຽງເໜືອ",
                        "contactInfo": "020 56681336 - 020 56488886"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280527",
                    "trackingPlatform": "8260506136004",
                    "trackingId": "8260506136004",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Nalee",
                    "receiverPhone": "96290686",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:39.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:51:55.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງນ້ຳທາ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານຂອນ",
                        "branch_address": "ບ້ານ ຂອນ",
                        "contactInfo": "030 9504106,  020 99635000"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280434",
                    "trackingPlatform": "8260508436144",
                    "trackingId": "8260508436144",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Sit taebgphachan",
                    "receiverPhone": "57879133",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:45.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:51:08.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ນາຊາຍ1(ນາຊາຍທອງ)",
                        "branch_address": "ບ້ານ ນາຊາຍ1",
                        "contactInfo": "ຂ້າງໂຮງຮຽນ ມສ ນາຊາຍ,020 93 178 404 "
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280323",
                    "trackingPlatform": "8260509983988",
                    "trackingId": "8260509983988",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "somthiap",
                    "receiverPhone": "98368600",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:30.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:50:13.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ສາລະວັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເມືອງລະຄອນເພັງ",
                        "branch_address": "ບ້ານ ລະຄອນເພັງ",
                        "contactInfo": "99382631"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280172",
                    "trackingPlatform": "8260508789948",
                    "trackingId": "8260508789948",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "seuk koemany",
                    "receiverPhone": "58981386",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:26.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:48:56.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫ້ວຍຢໍ້ ( ເຂື່ອນນໍ້າຂອງ )",
                        "branch_address": "ບ້ານ ຫ້ວຍຢໍ້",
                        "contactInfo": "ສາຂາ ຫ້ວຍຢໍ້ ເຂື່ອນໄຟຟ້າ ນ້ຳຂອງເບີໂທ: 020 23935639,  020 58169019"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296280077",
                    "trackingPlatform": "8260504883327",
                    "trackingId": "8260504883327",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "somsakoun",
                    "receiverPhone": "97531628",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:18.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:48:15.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຊະນະໄຊ (ເມືອງ ປາກຊັນ)",
                        "branch_address": "ບ້ານ ຊະນະໄຊ",
                        "contactInfo": "96332462,0309924246,99761105 ຕິດຕໍ່ພະນັກງານ(95323661  ເຈົ້າຂອງສາຂາ)"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279957",
                    "trackingPlatform": "8260509934576",
                    "trackingId": "8260509934576",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Loung junthongkham",
                    "receiverPhone": "56428919",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:22.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:47:21.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຄຳມ່ວນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ເລົ່າໂພໄຊ(ທ່າແຂກ)",
                        "branch_address": "ບ້ານ ເລົ່າໂພໄຊ",
                        "contactInfo": "  77786718 / 98185718"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279777",
                    "trackingPlatform": "8260506170616",
                    "trackingId": "8260506170616",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ທ ແສງນະພາ",
                    "receiverPhone": "55399417",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:35.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:46:05.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຄຳມ່ວນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບໍໂພນຕີ້ວ(ຄຳມ່ວນ)",
                        "branch_address": "ບ້ານ ໂພນບໍຕີ້ວ",
                        "contactInfo": "02056695570 02076270002"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279697",
                    "trackingPlatform": "8260505638096",
                    "trackingId": "8260505638096",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ນິກອນ ໄມໜໍແກ້ວ",
                    "receiverPhone": "23432509",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:42:52.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:45:30.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຄຳມ່ວນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໜອງບົກ(ເມືອງ ໜອງບົກ)",
                        "branch_address": "ບ້ານ ໜອງບົກ",
                        "contactInfo": "ປະສານວຽກ 02059659969.59524535 | 0309305394"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279581",
                    "trackingPlatform": "8260509895919",
                    "trackingId": "8260509895919",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Keo  khamkeo",
                    "receiverPhone": "98411183",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:15.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:44:36.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານພູໝອກ(ໂຮງໝໍແຂວງຫຼວງພະບາງ)",
                        "branch_address": "ບ້ານ ພູໝອກ",
                        "contactInfo": "020 97042119/  ໃກ້ກັບໂຮງໝໍແຂວງຫຼວງພະບາງ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279441",
                    "trackingPlatform": "8260508264980",
                    "trackingId": "8260508264980",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Lattana xaiyachuk",
                    "receiverPhone": "95897092",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:11.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:43:32.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ນາໄຊ",
                        "branch_address": "ບ້ານ ນາໄຊ",
                        "contactInfo": "ບ້ານ ນາໄຊ ເມືອງ ໄຊເສດຖາ ນະຄອນຫຼວງວຽງ.ຂ້າງຕູ້ATM JDB.ຮ້ານthee koff ເສັ້ນຕະຫຼາດນ້ອຍນາໄຊ020 97 758 993"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279353",
                    "trackingPlatform": "8260501341457",
                    "trackingId": "8260501341457",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Phim tnx",
                    "receiverPhone": "57548262",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:46:01.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:42:48.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຊຽງຂວາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຍອດລ້ຽງ (ເມືອງ ຄຳ)",
                        "branch_address": "ຍອດລ້ຽງ",
                        "contactInfo": "02092286998"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279257",
                    "trackingPlatform": "8260508034929",
                    "trackingId": "8260508034929",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Thoungxai",
                    "receiverPhone": "96995004",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:41:54.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:41:57.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຫຼັກ15(ໄຊທານີ)",
                        "branch_address": "ບ້ານ ດົງສ້າງຫີນ",
                        "contactInfo": "020 94014442"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296279141",
                    "trackingPlatform": "8260501979884",
                    "trackingId": "8260501979884",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "singtongviphone",
                    "receiverPhone": "91501366",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:53.000Z",
                    "width": 60,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 15000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:40:59.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຜາສັງ(ເມືອງເຟືອງ)",
                        "branch_address": "ບ້ານ ຜາສັງ",
                        "contactInfo": "ໂທ & WhatsApp : 76523871"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296278476",
                    "trackingPlatform": "8260508121264",
                    "trackingId": "8260508121264",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 139000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Toutou",
                    "receiverPhone": "54492249",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:45:58.000Z",
                    "width": 55,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 14000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:35:58.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໜອງບົວທອງ(ເມືອງ ສີໂຄດຕະບອງ)",
                        "branch_address": "ບ້ານ ໜອງບົວທອງ",
                        "contactInfo": "77806272 77780806 99194149 ພິກັດ: ສາມແຍກໜອງບົວທອງໃຕ້ ຕິດກັບຮ້ານລ້າງອັດສິດລົດສິດ"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296278311",
                    "trackingPlatform": "8260501487861",
                    "trackingId": "8260501487861",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Ari noystep",
                    "receiverPhone": "97474335",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:44:00.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:34:50.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຈຳປາສັກ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ບ້ານລົມສັກເໜືອ(ເມືອງ ບາຈຽງ)",
                        "branch_address": "ບ້ານ ລົມສັກເໜືອ",
                        "contactInfo": "ຮ້ານຕັ້ງຢູ່ຫຼັກ12ທາງໄປປະທຸມພອນ ໂທແລະແອັບ ເບີແອັບ02094315653/ໂທ93147480"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "296278126",
                    "trackingPlatform": "8260508614852",
                    "trackingId": "8260508614852",
                    "itemName": "ເຄື່ອງໃຊ້",
                    "itemValueKIP": 169000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Sack sengaloun",
                    "receiverPhone": "96789543",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "TRANSIT_TO_DEST_BRANCH",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-19T14:43:01.000Z",
                    "width": 50,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 13000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-19T13:33:40.000Z",
                    "destReceiveDate": null,
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຫຼວງພະບາງ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ສີມຸງຄຸນ(ເມືອງນ່ານ)",
                        "branch_address": "ບ້ານ ສີມຸງຄຸນ",
                        "contactInfo": "96591120"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "originReceiveBy": {
                        "first_name": "ທ້າວ ຄອນສະຫວັນ",
                        "phone_number": 58765428
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "295908198",
                    "trackingPlatform": "8260481778705",
                    "trackingId": "8260481778705",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 99000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 0,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "H\"Hom kamkeo",
                    "receiverPhone": "91155131",
                    "description": "",
                    "isSummary": 0,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "DEST_BRANCH_RECEIVED_FORWARD",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-18T03:25:45.000Z",
                    "width": 40,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 8000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-17T10:43:38.000Z",
                    "destReceiveDate": "2026-02-19T10:17:30.000Z",
                    "sendCompleteDate": null,
                    "isBackward": 0,
                    "billNumber": 0,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ຄຳມ່ວນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຍົມມະລາດ",
                        "branch_address": "ບ້ານ ຍົມມະລາດ",
                        "contactInfo": "55570739 "
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "294409965",
                    "trackingPlatform": "8260405324680",
                    "trackingId": "8260405324680",
                    "itemName": "ສຜ",
                    "itemValueKIP": 191000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 191000,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "jik visando",
                    "receiverPhone": "52570629",
                    "description": "",
                    "isSummary": 1,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "COMPLETED",
                    "contactStatus": "CONTACTED",
                    "originSendDate": "2026-02-09T07:42:11.000Z",
                    "width": 40,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 8000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-09T06:07:43.000Z",
                    "destReceiveDate": "2026-02-11T02:52:11.000Z",
                    "sendCompleteDate": "2026-02-12T09:19:06.000Z",
                    "isBackward": 0,
                    "billNumber": 57543182,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອັດຕະປື"
                    },
                    "originBranchId": {
                        "id_branch": "843",
                        "branch_name": "ສາຂາ ໂຊກຄຳ(ເມືອງ ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ທ່າຫິນ(ອັດຕະປື)",
                        "branch_address": "ບ້ານ ທ່າຫິນ",
                        "contactInfo": "58139393"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ນາງ ຈັນທະລາ ",
                        "phone_number": 58862370
                    },
                    "originReceiveBy": {
                        "first_name": "ນາງ ຈັນທະລາ ",
                        "phone_number": 58862370
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "294316426",
                    "trackingPlatform": "8260390672476",
                    "trackingId": "8260390672476",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 119000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 119000,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "ນິດ ຄັນທະລາ",
                    "receiverPhone": "58881432",
                    "description": "",
                    "isSummary": 1,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "COMPLETED",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-08T10:21:43.000Z",
                    "width": 40,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 8000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-08T10:06:23.000Z",
                    "destReceiveDate": "2026-02-11T04:18:38.000Z",
                    "sendCompleteDate": "2026-02-13T06:15:55.000Z",
                    "isBackward": 0,
                    "billNumber": 57553173,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ອຸດົມໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ໂຮມໄຊ(ເມືອງນາໝໍ້)",
                        "branch_address": "ບ້ານ ໂຮມໄຊ",
                        "contactInfo": "ເບີຕິດຕໍ່ສາຂາ 02093030168 ເວລາເປີດບໍລິການ8:30-17:30"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "294316207",
                    "trackingPlatform": "8260394582528",
                    "trackingId": "8260394582528",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 138000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 138000,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Had Had",
                    "receiverPhone": "78087014",
                    "description": "",
                    "isSummary": 1,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "COMPLETED",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-08T10:25:41.000Z",
                    "width": 40,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 8000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-08T10:05:47.000Z",
                    "destReceiveDate": "2026-02-09T07:48:43.000Z",
                    "sendCompleteDate": "2026-02-10T07:40:19.000Z",
                    "isBackward": 0,
                    "billNumber": 57523359,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ແຂວງ ບໍລິຄຳໄຊ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ທາງແບ່ງຫລັກ 20 ",
                        "branch_address": "ບ້ານ ວຽງຄໍາ",
                        "contactInfo": "ໂທ: 030 443 6851, ແອັບ: 020 7663 5154"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                },
                {
                    "_id": "294315892",
                    "trackingPlatform": "8260391937986",
                    "trackingId": "8260391937986",
                    "itemName": "ຄຊ",
                    "itemValueKIP": 211000,
                    "itemValueTHB": 0,
                    "itemValueUSD": 0,
                    "realItemValueKIP": 211000,
                    "realItemValueTHB": 0,
                    "realItemValueUSD": 0,
                    "receiverName": "Toutou",
                    "receiverPhone": "59999993",
                    "description": "",
                    "isSummary": 1,
                    "destSendDate": null,
                    "charge_on_shop": 0,
                    "itemStatus": "COMPLETED",
                    "contactStatus": "NOT_CONTACT",
                    "originSendDate": "2026-02-08T10:21:35.000Z",
                    "width": 45,
                    "weight": 1,
                    "isCod": "1",
                    "isExtraItem": 0,
                    "packagePrice": 12000,
                    "isInsurance": 0,
                    "insuranceAmount": 0,
                    "originReceiveDate": "2026-02-08T10:04:55.000Z",
                    "destReceiveDate": "2026-02-09T05:18:07.000Z",
                    "sendCompleteDate": "2026-02-09T11:16:02.000Z",
                    "isBackward": 0,
                    "billNumber": 57513272,
                    "fcm_token_staff": null,
                    "status_pay_by_bank": null,
                    "originProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "destProvinceId": {
                        "provinceName": "ນະຄອນຫຼວງວຽງຈັນ"
                    },
                    "originBranchId": {
                        "id_branch": "84",
                        "branch_name": "ສາຂາ ວຽງຈະເລີນ(ໄຊເສດຖາ)",
                        "mainBranches": 0,
                        "merchant_id_kip": null,
                        "merchant_id_thb": null,
                        "merchant_id_usd": null
                    },
                    "destBranchId": {
                        "branch_name": "ສາຂາ ຈິນາຍໂມ້(ເມືອງ ສີສັດຕະນາກ)",
                        "branch_address": "ບ້ານ ຈອມເພັດໃຕ້",
                        "contactInfo": "ຕິດຕໍ່ຝາກເຄື່ອງ-ນຳບີນ-ຮັບເຄື່ອງ 02079993322 ແລະ 02098828340"
                    },
                    "customerId": {
                        "id_list": "439646",
                        "full_name": "ຕີນາ ທິບພະວົງ",
                        "contact_info": "76681339"
                    },
                    "createdBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "originReceiveBy": {
                        "first_name": "ສາຂາ ",
                        "phone_number": 55966379
                    },
                    "providedBy": {
                        "_id": "1"
                    }
                }
            ]
        }
    }
}`

	var response custommodel.ItemsV2Response
	if err := json.Unmarshal([]byte(rawJSON), &response); err != nil {
		log.Println("Error parsing hardcoded JSON:", err)
	}
	return response
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
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set your headers to match the real browser request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en,th;q=0.9,en-US;q=0.8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("Origin", "https://app.anousith.express")
	req.Header.Set("Referer", "https://app.anousith.express/")

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
