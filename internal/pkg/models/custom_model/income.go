package custommodel

import "time"

// FinancialStatus Struct สำหรับรองรับผลลัพธ์จาก SQL Dashboard
type FinancialStatus struct {
	TotalIncome  float64 `json:"total_income" gorm:"column:total_income"`
	TotalExpense float64 `json:"total_expense" gorm:"column:total_expense"`
	NetBalance   float64 `json:"net_balance" gorm:"column:net_balance"`
}

// StatementItem ใช้รับผลลัพธ์จาก SQL Union
type StatementItem struct {
	TransactionDate time.Time `json:"transaction_date" gorm:"column:transaction_date"`
	Type            string    `json:"type" gorm:"column:type"` // INCOME หรือ EXPENSE
	CategoryName    string    `json:"category_name" gorm:"column:category_name"`
	Description     string    `json:"description" gorm:"column:description"`
	Amount          float64   `json:"amount" gorm:"column:amount"`
	NetAmount       float64   `json:"net_amount" gorm:"column:net_amount"`
}
