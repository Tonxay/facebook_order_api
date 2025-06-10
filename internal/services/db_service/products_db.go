package dbservice

import (
	"context"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/query"

	"gorm.io/gorm"
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

func GetProducts(db *gorm.DB) ([]custommodel.Products, error) {

	var products []custommodel.Products

	tx := db.Table(models.TableNameProduct + " p")

	tx = tx.Select("p.id,p.name")
	tx = tx.Preload("Promotions", func(db *gorm.DB) *gorm.DB {
		tx := db.Where("status = ?", "active")
		return tx
	})

	tx = tx.Preload("ProductDetails", func(db *gorm.DB) *gorm.DB {

		tx := db.Where("status = ?", "active")

		tx = tx.Preload("Sizes", func(db *gorm.DB) *gorm.DB {
			tx := db.
				Select(
					`    sizes.id,
					     sizes.size,
					     sizes.price,
					     sizes.product_detail_id,
					     SUM(s.remaining) AS total_remaining
						  
				    `,
				)

			tx = tx.Joins("LEFT JOIN " + models.TableNameStockProductDetail + " s ON s.size_id = sizes.id")

			tx.Where("s.remaining > ? AND s.status = ?", 0, "active")

			tx = tx.Group("sizes.id, sizes.size,sizes.product_detail_id")

			return tx
		})
		return tx
	})

	err := tx.Find(&products).Error

	return products, err
}
