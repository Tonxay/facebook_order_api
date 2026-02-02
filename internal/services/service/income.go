package service

import (
	"go-api/internal/config/middleware"
	"go-api/internal/config/presenters"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	dbservice "go-api/internal/services/db_service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateIncomesInBatch(c *fiber.Ctx) error {
	// 1. เปลี่ยนตัวแปรรับค่าเป็น Slice ของ IncomeInput (เพราะรับแบบ Batch)
	var inputs []models.Income

	// ใช้ Middleware เดิมของคุณ (ตรวจสอบให้แน่ใจว่า Middleware นี้รองรับ Slice/Array)
	if err := middleware.ParseAndValidateBody(c, &inputs); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// 3. เรียกใช้ DB Service ที่เราสร้างไว้
	db := gormpkg.GetDB()
	if err := dbservice.CreateIncomesInBatch(db, inputs); err != nil {
		// Log Error ไว้ดูเองด้วยจะดีมาก
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create incomes: "+err.Error())
	}

	// 4. ส่งค่ากลับ (201 Created)
	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(fiber.Map{
		"message": "success",
		"count":   len(inputs),
		"data":    inputs, // ส่งข้อมูลที่มี ID ที่ Gen จาก DB กลับไป
	}))
}

// GetIncomeCategoriesHandler: API สำหรับดึงหมวดหมู่รายรับ
// GET /api/categories/income
func GetIncomeCategoriesHandler(c *fiber.Ctx) error {
	db := gormpkg.GetDB()

	categories, err := dbservice.GetIncomeCategories(db)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลหมวดหมู่ได้")
	}

	return c.JSON(presenters.ResponseSuccess(categories))
}

// CreateIncomeCategoryHandler: API สำหรับเพิ่มหมวดหมู่ใหม่
// POST /api/categories/income
func CreateIncomeCategoryHandler(c *fiber.Ctx) error {
	var input models.IncomeCategory

	// รับค่า Body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Validate เบื้องต้น (ถ้าชื่อว่างเปล่า)
	if input.CategoryName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Category Name is required"})
	}

	db := gormpkg.GetDB()
	if err := dbservice.CreateIncomeCategory(db, &input); err != nil {
		// เช็คกรณีชื่อซ้ำ (Unique Constraint)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถสร้างหมวดหมู่ได้ (ชื่ออาจซ้ำ)"})
	}

	return c.Status(fiber.StatusCreated).JSON(presenters.ResponseSuccess(input))
}

func GetIncomesHandler(c *fiber.Ctx) error {
	// 1. รับค่า Query Params (ถ้าไม่ส่งมา ให้ Default เป็นเดือน/ปี ปัจจุบัน)
	now := time.Now()

	month := c.QueryInt("month", int(now.Month())) // Default: เดือนปัจจุบัน
	year := c.QueryInt("year", now.Year())         // Default: ปีปัจจุบัน

	// Validate ค่าเดือนนิดหน่อย
	if month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid month (must be 1-12)",
		})
	}

	// 2. เรียก DB Service
	db := gormpkg.GetDB()
	incomes, err := dbservice.GetIncomesByMonth(db, month, year)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch incomes",
		})
	}

	// 3. ส่ง Response กลับแบบมาตรฐาน
	return c.Status(fiber.StatusOK).JSON(presenters.ResponseSuccess(fiber.Map{
		"month":   month,
		"year":    year,
		"count":   len(incomes),
		"results": incomes,
	}))
}

func GetFinancialStatusHandler(c *fiber.Ctx) error {
	// 1. รับค่า date จาก Query Param (เช่น "2024-01-01")
	targetDate := c.Query("date")

	// 2. ถ้าไม่ได้ส่งมา (เป็นค่าว่าง) ให้ใช้วันปัจจุบันเป็นค่า Default
	if targetDate == "" {
		targetDate = time.Now().Format("2006-01-02")
	}

	// 3. เรียก DB Service (ตัวใหม่ที่รับค่า date)
	db := gormpkg.GetDB()
	status, err := dbservice.GetFinancialStatusByDate(db, targetDate)

	if err != nil {
		// Log error ไว้หน่อยก็ดีครับ
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch financial status")
	}

	// 4. ส่งผลลัพธ์กลับแบบมาตรฐาน
	return c.Status(fiber.StatusOK).JSON(presenters.ResponseSuccess(status))
}

func GetStatementHandler(c *fiber.Ctx) error {
	// 1. รับค่า Params (ถ้าไม่ส่งมา ให้ตั้งค่า Default)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	// Default วันที่: วันแรกของเดือน ถึง วันนี้
	now := time.Now()
	startDate := c.Query("start_date", time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
	endDate := c.Query("end_date", now.Format("2006-01-02"))

	// 2. เรียก DB Service
	db := gormpkg.GetDB()
	items, err := dbservice.GetStatement(db, startDate, endDate, page, limit)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch statement")
	}

	// 3. ส่งผลลัพธ์กลับ
	return c.Status(fiber.StatusOK).JSON(presenters.ResponseSuccess(fiber.Map{
		"page":       page,
		"limit":      limit,
		"start_date": startDate,
		"end_date":   endDate,
		"count":      len(items),
		"data":       items,
	}))
}
