package dbservice

import (
	"context"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/query"
)

func CreateStockProductDetail(stockProductDetail *models.StockProductDetail, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.StockProductDetail
	err := daq.WithContext(ctx).Create(stockProductDetail)
	return err
}

func CreateStockDetails(stockDetail *models.StockDetail, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.StockDetail
	err := daq.WithContext(ctx).Create(stockDetail)
	return err
}

func UpdateStockProductDetail(id string, remaining int32, status string, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.StockProductDetail
	_, err := daq.WithContext(ctx).Where(daq.ID.Eq(id)).UpdateSimple(daq.Remaining.Value(remaining), daq.Status.Value(status))
	return err
}
