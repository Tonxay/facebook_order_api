// in file: dbservice/expense_service.go

package dbservice

import (
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"math"

	"gorm.io/gorm"
)

func GetFullExpenseDetails(db *gorm.DB, filter custommodel.ExpenseFilter) (*custommodel.PaginatedResponse, error) {
	var results []custommodel.FullExpenseDetail
	var totalRows int64

	// Set pagination defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	// Base query
	query := db.Table("expenses e").
		Joins("LEFT JOIN products p ON e.product_id = p.id").
		Joins("LEFT JOIN expense_categories ec ON e.category_id = ec.id").
		Joins("LEFT JOIN suppliers s ON e.supplier_id = s.id")

	// Apply filters conditionally
	if filter.Description != "" {
		query = query.Where("e.description ILIKE ?", "%"+filter.Description+"%")
	}
	if filter.CategoryID != "" {
		query = query.Where("e.category_id = ?", filter.CategoryID)
	}
	if filter.SupplierID != "" {
		query = query.Where("e.supplier_id = ?", filter.SupplierID)
	}

	// **Add date range filter logic**
	if filter.StartDate != "" {
		query = query.Where("e.expense_date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("e.expense_date <= ?", filter.EndDate)
	}

	// 1. Get the total count
	if err := query.Model(&custommodel.FullExpenseDetail{}).Count(&totalRows).Error; err != nil {
		return nil, err
	}

	// 2. Apply pagination and fetch the data
	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Select(`
            e.*,
            p.name AS product_name,
            ec.category_name AS category_name,
            s.supplier_name AS supplier_name
        `).
		Order("e.expense_date DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 3. Assemble the paginated response
	response := &custommodel.PaginatedResponse{
		Data:       results,
		TotalRows:  totalRows,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: int64(math.Ceil(float64(totalRows) / float64(filter.PageSize))),
	}

	return response, nil
}

// CreateExpensesInBatch creates multiple expense records in a single database transaction.
func CreateExpenses(db *gorm.DB, expenses []*custommodel.ExpenseInput) ([]*models.Expense, error) {

	// Convert input models to database models
	var expenseModels []*models.Expense
	for _, input := range expenses {
		expenseModels = append(expenseModels, &models.Expense{
			ExpenseDate: input.ExpenseDate,
			Amount:      input.Amount,
			Description: *input.Description,
			ReceiptURL:  *input.ReceiptURL,
			CategoryID:  input.CategoryID,
			SupplierID:  *input.SupplierID,
			ProductID:   *input.ProductID,
		})
	}

	// Use a transaction to ensure all records are created or none are.
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&expenseModels).Error; err != nil {
			// If any record fails, the transaction is rolled back.
			return err
		}
		// If all is well, the transaction is committed.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return expenseModels, nil
}

// CreateExpensesInBatch creates multiple expense records from a slice of inputs.
func CreateExpensesInBatch(db *gorm.DB, inputs []*custommodel.ExpenseInput) ([]*custommodel.Expense, error) {
	var expenseModels []*custommodel.Expense
	for _, input := range inputs {
		expenseModels = append(expenseModels, &custommodel.Expense{
			ExpenseDate: input.ExpenseDate,
			Amount:      input.Amount,
			Description: input.Description,
			ReceiptURL:  input.ReceiptURL,
			CategoryID:  input.CategoryID,
			SupplierID:  input.SupplierID,
			ProductID:   input.ProductID, // Maps the product ID if provided
		})
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&expenseModels).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return expenseModels, nil
}

func GetFullExpenseCategory(db *gorm.DB) ([]models.ExpenseCategory, error) {
	var results []models.ExpenseCategory

	err := db.Model(&results).Find(&results).Error

	return results, err
}

func GetFullExpenseSuppliers(db *gorm.DB) ([]models.Supplier, error) {
	var results []models.Supplier

	err := db.Model(&results).Find(&results).Error

	return results, err
}

func GetFullExpenseProducts(db *gorm.DB) ([]models.Product, error) {
	var results []models.Product

	err := db.Model(&results).Find(&results).Error

	return results, err
}
