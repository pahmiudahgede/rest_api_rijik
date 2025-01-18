package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"gorm.io/gorm"
)

func GetProductsByStoreID(storeID string, limit, offset int) ([]domain.Product, error) {
	var products []domain.Product
	query := config.DB.Preload("ProductImages").Preload("TrashDetail").Where("store_id = ?", storeID)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

func GetProductsByUserID(userID string, limit, offset int) ([]domain.Product, error) {
	var products []domain.Product
	query := config.DB.Preload("ProductImages").Preload("TrashDetail").Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

func GetProductByIDAndStoreID(productID, storeID string) (domain.Product, error) {
	var product domain.Product
	err := config.DB.Preload("ProductImages").Preload("TrashDetail").
		Where("id = ? AND store_id = ?", productID, storeID).
		First(&product).Error

	return product, err
}

func GetProductByID(productID string) (domain.Product, error) {
	var product domain.Product
	err := config.DB.Preload("ProductImages").Preload("TrashDetail").
		Where("id = ?", productID).First(&product).Error
	return product, err
}

func IsValidStoreID(storeID string) bool {
	var count int64
	err := config.DB.Model(&domain.Store{}).Where("id = ?", storeID).Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CreateProduct(product *domain.Product, images []domain.ProductImage) error {

	return config.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(product).Error; err != nil {
			return err
		}

		if len(images) > 0 {
			for i := range images {
				images[i].ProductID = product.ID
			}

			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateProduct(product *domain.Product, images []domain.ProductImage) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(product).Error; err != nil {
			return err
		}

		if len(images) > 0 {
			for i := range images {
				images[i].ProductID = product.ID
			}

			if err := tx.Where("product_id = ?", product.ID).Delete(&domain.ProductImage{}).Error; err != nil {
				return err
			}

			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func DeleteProduct(productID string) error {

	return config.DB.Where("id = ?", productID).Delete(&domain.Product{}).Error
}
