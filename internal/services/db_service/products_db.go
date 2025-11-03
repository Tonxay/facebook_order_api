package dbservice

import (
	"context"
	"encoding/json"
	"fmt"
	gormpkg "go-api/internal/pkg"
	"go-api/internal/pkg/models"
	custommodel "go-api/internal/pkg/models/custom_model"
	"go-api/internal/pkg/query"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

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

	tx = tx.Where("p.status  = ?", "active")
	tx = tx.Preload("Promotions", func(db *gorm.DB) *gorm.DB {
		tx := db.Where("status = ?", "active").Order("quentity ASC")
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

			// tx.Where("s.remaining >= ? AND s.status = ? OR s.status = ?", 0, "active", "out_stock")

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

	// tx = tx.Where("sd.status  = ?", "active")

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

func GetProductsResearch(db *gorm.DB) ([]models.ProductsResearch, error) {

	var products []models.ProductsResearch

	err := db.Model(&products).Find(&products).Error
	return products, err
}

// Struct to partially decode the response to check for the error - still useful
type LazadaPartialResponse struct {
	Ret []string `json:"ret"`
}

// Function to fetch product research data from Lazada and return parsed JSON
func GetProductsResearchFormLazada(additionalParams map[string]string) (map[string]interface{}, error) {
	baseURL := "https://www.lazada.co.th/catalog/"

	// Prepare query parameters
	queryParams := url.Values{}
	queryParams.Set("ajax", "true")

	if additionalParams != nil {
		for key, value := range additionalParams {
			queryParams.Set(key, value)
		}
	}

	fullURL := baseURL + "?" + queryParams.Encode()
	fmt.Println("Requesting URL:", fullURL)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Cookie", "__wpkreporterwid_=50b759dc-6a7c-45ec-3e7f-1bf09e6869e5; t_fv=1761202379012; t_uid=4gQAXa7TqRRCwPnqRZ8u8gwbtpDNpJyx; lzd_cid=ecb05c8d-e4bf-4f74-fa7d-246354971773; lzd_sid=1513ecdb18c828ddfdb8e038671b5b06; cna=y76AIZxHNX8CAYsFny+PucQY; lwrid=AgGaD9eari4LK5RmeOttxe5uI8qS; _gcl_au=1.1.1568506510.1761202380; xlly_s=1; lzd_uid=100093350650; sgcookie=E100ShaIPPlwZK9hOom1SzcEu2avZQxvsiUn2OgZdjy4G%2FBEPvoKKmFu5JzJLnHHjCWHfPv6Eem1hgEDkdLrEWMJmw6DRYkRMmY%2FEX4p2cJkIfA%3D; lzd_login_lastlogintype=GOOGLE; __itrace_wid=7f7d183d-eaf2-4758-beaf-095303e02637; _fbp=fb.2.1761203995112.342354004270367852; AMCV_126E248D54200F960A4C98C6%40AdobeOrg=-1124106680%7CMCIDTS%7C20385%7CMCMID%7C07278010371852321831629851061915852289%7CMCAAMLH-1761808795%7C3%7CMCAAMB-1761808795%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1761211195s%7CNONE%7CvVersion%7C5.2.0; hng=TH|th|THB|764; userLanguageML=th; _uetvid=aca3eaf0afe011f0a11c1fc3237b5bfc; isg=BJmZldqHFfP1xMll-pwJKsZLqINzJo3YTuakLbtDDUSowoFUBnVOqLmZxZ60rSUQ; t_sid=dS37MJ7B5a4jfDiXmNmDGf1JvzLIhGDR; utm_channel=NA; lzd_uti=%7B%22fpd%22%3A%229999-99-99%22%2C%22lpd%22%3A%229999-99-99%22%2C%22cnt%22%3A%220%22%7D; _tb_token_=eee35eba33d5e; _m_h5_tk=3398e28c17922356f481c661938750a3_1761295187860; _m_h5_tk_enc=d427e869d4f157d3eda831811b4d6eb9; lwrtk=AAIEaPuPs82r5LYwNUmpuaLbi0uqFpot+nwvqxdUU1f9yp45hG4IaFU=; epssw=10*_Bpss6C9hBO5kWxaORe0qfFFoRPcZs04FX8pQnclO-Jojp4It-ZcA6seU7V0jNFe1VV0U-XQO-J0U-upChHVCYiVmUsBFgaMZW_2krWoubm2ssssCpanXRWbtQS2-ocX4HXfzFsazk3kOOHnN_0gPZLMkC5hI-PJOeYdM38wDBp7sj5Ry9LBpPJ8weAkvLVPm34MStC5HiQ1jTZEOaXtO7D8Ietcg6snoEjSUa5n2tLZUWrbbRssssduFOss3hu_HXFl3ReWUAL6jsve-baQoiVqrRm9ddiv8nKXxcLaoiGsvJ4z_mNvMT6GNQtgm6L-; tfstk=gF-jBbYOVCKrG48dlGkzRtNDiGs_1YoU5R69KdE4BiIY65pdaGk07qX1f_OGgNCwbQ15m384bslcfddFfXlEYDJDFistTXygnsLFXOnN647xy1SX-XlEY0JA1wqZTCPy6K7RI_BYDrExFzB1dGBOkOH5wO63WSdOXYM5C9PA6GITw0BCCGC96GHWe_XOXtdOXY9RZOLzPJ6slsvjKyawbt2IvgC7XlK5exfphg6lE36fl6KcxlB9Vt_fvKtChSiFCeKCoBhzusLB7B6eg0EXOKdXJwKYwbsMUeO1dnHLR1ODen7JDAzhEGYvJiKsOSpfR3KVmUGQY6LkBh_JnXahUFpe0FOqZlC2RLd583PzYgKXGnQ5vgPbY66wTPw5-l65TYM7SPf_jsd3OAxdZZBlhakSFSLgtlE69hk7WnQAETuxFYNvS")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing request: %v\n", err)
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	log.Printf("Lazada Response Status: %s\n", resp.Status)

	// --- Decode the JSON into a map ---
	var responseData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		// If decoding fails, it might be HTML or invalid JSON
		log.Printf("Error decoding response body as JSON: %v. Body was: %s\n", err, string(bodyBytes))
		// Attempt to check if it's the specific blocking error even if full decode failed
		var partialResp LazadaPartialResponse
		if json.Unmarshal(bodyBytes, &partialResp) == nil && len(partialResp.Ret) > 0 && partialResp.Ret[0] == "FAIL_SYS_USER_VALIDATE" {
			log.Printf("Lazada blocked the request (CAPTCHA likely required, partial decode): %s\n", string(bodyBytes))
			// Return the partially decoded data (or nil) along with the specific error
			return responseData, fmt.Errorf("lazada request blocked (FAIL_SYS_USER_VALIDATE)")
		}
		// Return a general JSON parsing error
		return nil, fmt.Errorf("error decoding response body as JSON: %w", err)
	}

	// --- Check for the specific Lazada error *within* the decoded JSON ---
	if retVal, ok := responseData["ret"].([]interface{}); ok && len(retVal) > 0 {
		if retStr, ok := retVal[0].(string); ok && retStr == "FAIL_SYS_USER_VALIDATE" {
			log.Printf("Lazada blocked the request (CAPTCHA likely required): %v\n", responseData)
			// Return the decoded map along with the specific blocking error
			return responseData, fmt.Errorf("lazada request blocked (FAIL_SYS_USER_VALIDATE)")
		}
	}

	// If no block detected, return the decoded map
	log.Println("Lazada request successful and decoded.")
	return responseData, nil
}
