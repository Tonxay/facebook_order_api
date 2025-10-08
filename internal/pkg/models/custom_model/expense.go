package custommodel

import (
	"go-api/internal/pkg/models"
	"time"
)

// ExpenseFilter holds all possible query parameters for filtering and pagination.
type ExpenseFilter struct {
	Description string `query:"description"` // Search term for the description
	CategoryID  string `query:"category_id"` // Filter by a specific category UUID
	SupplierID  string `query:"supplier_id"` // Filter by a specific supplier UUID
	Page        int    `query:"page"`
	PageSize    int    `query:"page_size"`
	// Add date range fields
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// PaginatedResponse is a generic struct for API responses with pagination.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int64       `json:"total_pages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

// FullExpenseDetail holds the combined data from the expenses table and its joins.
type FullExpenseDetail struct {
	models.Expense         // Embed the original Expense model to get all its fields
	ProductName    *string `gorm:"column:product_name" json:"product_name"`
	CategoryName   string  `gorm:"column:category_name" json:"category_name"`
	SupplierName   *string `gorm:"column:supplier_name" json:"supplier_name"`
}

// ExpenseInput is used for creating and updating expenses.
type ExpenseInput struct {
	ExpenseDate time.Time `json:"expense_date" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Description *string   `json:"description"`
	ReceiptURL  *string   `json:"receipt_url"`
	CategoryID  string    `json:"category_id" binding:"required"`
	SupplierID  *string   `json:"supplier_id"`
	ProductID   *string   `json:"product_id"`
}

const TableNameExpense = "expenses"

// Expense mapped from table <expenses>
type Expense struct {
	ID          string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	ExpenseDate time.Time `gorm:"column:expense_date;not null" json:"expense_date"`
	Amount      float64   `gorm:"column:amount;not null" json:"amount"`
	Description *string   `gorm:"column:description" json:"description"`
	ReceiptURL  *string   `gorm:"column:receipt_url" json:"receipt_url"`
	ProductID   *string   `gorm:"column:product_id" json:"product_id"`
	CategoryID  string    `gorm:"column:category_id;not null" json:"category_id"`
	SupplierID  *string   `gorm:"column:supplier_id" json:"supplier_id"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName Expense's table name
func (*Expense) TableName() string {
	return TableNameExpense
}

type ExpenseSummary struct {
	CategoryID   string  `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	TotalAmount  float64 `json:"total_amount" db:"total_amount"`
}
