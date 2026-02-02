package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"time"

	"gorm.io/gorm"
)

// CreateIncomesInBatch รับข้อมูลเป็น struct ของ Income โดยตรง แล้วบันทึกทีละหลายรายการ
func CreateIncomesInBatch(db *gorm.DB, incomes []models.Income) error {
	// เช็คก่อนว่ามีข้อมูลไหม ถ้าไม่มีก็ไม่ต้องทำอะไร
	if len(incomes) == 0 {
		return nil
	}

	// ใช้ CreateInBatches ของ GORM
	// พารามิเตอร์ 100 คือขนาดของ Batch (บันทึกทีละ 100 รายการ) เพื่อไม่ให้ SQL query ยาวเกินไป
	result := db.CreateInBatches(&incomes, 100)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetIncomeCategories: ดึงรายชื่อหมวดหมู่ทั้งหมด (สำหรับ Dropdown)
func GetIncomeCategories(db *gorm.DB) ([]models.IncomeCategory, error) {
	var categories []models.IncomeCategory

	// ดึงข้อมูลทั้งหมด เรียงตามชื่อ
	result := db.Order("category_name ASC").Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// CreateIncomeCategory: สร้างหมวดหมู่ใหม่
func CreateIncomeCategory(db *gorm.DB, category *models.IncomeCategory) error {
	return db.Create(category).Error
}

// GetIncomesByMonth ดึงรายรับประจำเดือนที่ระบุ
func GetIncomesByMonth(db *gorm.DB, month int, year int) ([]models.Income, error) {
	var incomes []models.Income

	// 1. สร้างช่วงเวลา StartDate และ EndDate ของเดือนนั้น
	// เช่น month=1, year=2024 -> 2024-01-01 ถึง 2024-02-01
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0) // บวกไป 1 เดือนเพื่อเป็นจุดตัด

	// 2. Query โดยใช้ช่วงเวลา (Range Query) จะเร็วกว่าการใช้ Function แปลงวันที่
	// เลือกเฉพาะรายการที่ >= วันแรกของเดือน และ < วันแรกของเดือนถัดไป
	err := db.Where("income_date >= ? AND income_date < ?", startDate, endDate).
		Order("income_date DESC"). // เรียงจากใหม่ไปเก่า
		Find(&incomes).Error

	if err != nil {
		return nil, err
	}

	return incomes, nil
}

// GetFinancialStatusByDate เลือกดูยอดของเดือนตามวันที่ที่ส่งเข้าไป (เช่น "2024-01-01")
func GetFinancialStatusByDate(db *gorm.DB, targetDate string) (*custommodel.FinancialStatus, error) {
	var result custommodel.FinancialStatus

	// ใช้ ? (Placeholder) เพื่อส่งตัวแปรเข้าไปแทน CURRENT_DATE
	query := `
		SELECT 
			COALESCE(SUM(i.amount), 0) AS total_income,
			COALESCE(SUM(e.amount), 0) AS total_expense,
			(COALESCE(SUM(i.amount), 0) - COALESCE(SUM(e.amount), 0)) AS net_balance
		FROM 
			(SELECT amount FROM incomes WHERE DATE_TRUNC('month', income_date) = DATE_TRUNC('month', ?::date)) i
		FULL OUTER JOIN 
			(SELECT amount FROM expenses WHERE DATE_TRUNC('month', expense_date) = DATE_TRUNC('month', ?::date)) e 
		ON 1=1;
	`

	// ส่ง targetDate เข้าไป 2 ที่ (สำหรับ incomes และ expenses)
	if err := db.Raw(query, targetDate, targetDate).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatement ดึงรายการเคลื่อนไหว (รวมรับ-จ่าย)
func GetStatement(db *gorm.DB, startDate, endDate string, page, limit int) ([]custommodel.StatementItem, error) {
	var results []custommodel.StatementItem

	// 1. SQL หลักของคุณ (Union)
	// ผมเอามาห่อด้วย subquery "all_transactions" เพื่อให้กรองวันที่ง่ายขึ้น
	sqlQuery := `
		SELECT * FROM (
			-- ดึงรายรับ
			SELECT 
				income_date as transaction_date,
				'INCOME' as type,
				ic.category_name,
				i.description,
				i.amount,
				i.amount as net_amount
			FROM incomes i
			JOIN income_categories ic ON i.category_id = ic.id

			UNION ALL

			-- ดึงรายจ่าย
			SELECT 
				expense_date as transaction_date,
				'EXPENSE' as type,
				ec.category_name,
				e.description,
				e.amount,
				-e.amount as net_amount
			FROM expenses e
			JOIN expense_categories ec ON e.category_id = ec.id
		) AS all_transactions
	`

	// 2. สร้างเงื่อนไข Filter และ Pagination
	// เงื่อนไข: WHERE transaction_date BETWEEN ? AND ?
	sqlQuery += " WHERE transaction_date BETWEEN ? AND ?"

	// เรียงลำดับ: ORDER BY transaction_date DESC
	sqlQuery += " ORDER BY transaction_date DESC, type"

	// Pagination: LIMIT ? OFFSET ?
	sqlQuery += " LIMIT ? OFFSET ?"

	// คำนวณ Offset
	offset := (page - 1) * limit

	// 3. รัน Query พร้อมส่ง Parameter เข้าไปตามลำดับ (?)
	// params: startDate, endDate, limit, offset
	if err := db.Raw(sqlQuery, startDate, endDate, limit, offset).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
