package dbservice

import (
	"context"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	"go-api/internal/pkg/query"
)

func CreateCategory(category *models.Category, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.Category
	err := daq.WithContext(ctx).Create(category)
	return err
}

func CreateProduct(product *models.Product, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.Product
	err := daq.WithContext(ctx).Create(product)
	return err
}

func CreateProductDetail(productDetail *models.ProductDetail, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.ProductDetail
	err := daq.WithContext(ctx).Create(productDetail)
	return err
}

func CreateProductSize(size *models.Size, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.Size
	err := daq.WithContext(ctx).Create(size)
	return err
}
