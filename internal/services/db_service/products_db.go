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

func CreateProduct(db *gorm.DB, product *models.Product, ctx context.Context) error {
	query.SetDefault(db)
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
func CreateProductDetailList(db *gorm.DB, productDetail []*models.ProductDetail, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.ProductDetail
	err := daq.WithContext(ctx).CreateInBatches(productDetail, 100)
	return err
}
func CreateProductSize(size *models.Size, ctx context.Context) error {
	query.SetDefault(gormpkg.GetDB())
	daq := query.Q.Size
	err := daq.WithContext(ctx).Create(size)
	return err
}

func CreateProductSizeList(db *gorm.DB, sizes []*models.Size, ctx context.Context) error {
	query.SetDefault(db)
	daq := query.Q.Size
	err := daq.WithContext(ctx).CreateInBatches(sizes, 100)
	return err
}
func GetProducts(db *gorm.DB) ([]custommodel.Products, error) {

	var products []custommodel.Products

	tx := db.Table(models.TableNameProduct + " p")

	tx = tx.Select("p.id,p.name,SUM(sd.remaining) AS total_amounts")

	tx = tx.Joins("LEFT JOIN " + models.TableNameProductDetail + " pd ON pd.product_id = p.id").
		Joins("LEFT JOIN " + models.TableNameStockProductDetail + " sd ON sd.product_detail_id = pd.id")

	tx = tx.Where("sd.status  = ?", "active")
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

			tx.Where("s.remaining >= ? AND s.status = ? OR s.status = ?", 0, "active", "out_stock")

			tx = tx.Group("sizes.id, sizes.size,sizes.product_detail_id")

			return tx
		})
		return tx
	})

	err := tx.Group(`p.id,p.name`).Find(&products).Error

	return products, err
}

func GetProductsForStock(db *gorm.DB) ([]custommodel.Products, error) {

	var products []custommodel.Products

	tx := db.Table(models.TableNameProduct + " p")

	tx = tx.Select("p.id,p.name,SUM(sd.remaining) AS total_amounts")

	tx = tx.Joins("LEFT JOIN " + models.TableNameProductDetail + " pd ON pd.product_id = p.id").
		Joins("LEFT JOIN " + models.TableNameStockProductDetail + " sd ON sd.product_detail_id = pd.id")

	tx = tx.Where("sd.status  = ?", "active")

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

			tx.Where("s.remaining >= ? AND s.status = ? OR s.status = ?", 0, "active", "out_stock")

			tx = tx.Group("sizes.id, sizes.size,sizes.product_detail_id")

			return tx
		})
		return tx
	})

	err := tx.Group(`p.id,p.name`).Find(&products).Error

	return products, err
}
