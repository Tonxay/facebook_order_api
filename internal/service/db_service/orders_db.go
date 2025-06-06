package dbservice

import (
	"context"
	"fmt"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/query"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *models.Order, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.Order
	err := daq.WithContext(ctx).Create(order)
	return err
}

func CreateOrderDetails(db *gorm.DB, orderDetail []*models.OrderDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.OrderDetail

	// Insert in batches, handle error
	err := daq.WithContext(ctx).CreateInBatches(orderDetail, 100)
	if err != nil {
		// Optional: Log or wrap for context
		return fmt.Errorf("failed to create order details: %w", err)
	}

	return nil
}
