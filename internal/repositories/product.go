package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"gorm.io/gorm"
)

func GetProductsByUserID(userID string, limit, offset int) ([]domain.Product, error) {
	var products []domain.Product
	query := config.DB.Preload("ProductImages").Preload("TrashDetail").Where("user_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

func GetProductByIDAndUserID(productID, userID string) (domain.Product, error) {
	var product domain.Product
	err := config.DB.Preload("ProductImages").Preload("TrashDetail").
		Where("id = ? AND user_id = ?", productID, userID).
		First(&product).Error

	return product, err
}

func GetProductByID(productID string) (domain.Product, error) {
	var product domain.Product
	err := config.DB.Preload("ProductImages").Preload("TrashDetail").
		Where("id = ?", productID).First(&product).Error
	return product, err
}

func CreateProduct(product *domain.Product, images []domain.ProductImage) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {

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
	return err
}

func UpdateProduct(product *domain.Product, images []domain.ProductImage) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&domain.Product{}).Where("id = ?", product.ID).Updates(product).Error; err != nil {
			return err
		}

		if err := tx.Where("product_id = ?", product.ID).Delete(&domain.ProductImage{}).Error; err != nil {
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

func DeleteProduct(productID, userID string) (int64, error) {
	result := config.DB.Where("id = ? AND user_id = ?", productID, userID).Delete(&domain.Product{})
	return result.RowsAffected, result.Error
}
