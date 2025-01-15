package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetAllProducts(limit, offset int) ([]domain.Product, error) {
	var products []domain.Product

	query := config.DB.Preload("ProductImages").Preload("TrashDetail")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductByID(productID string) (domain.Product, error) {
	var product domain.Product
	err := config.DB.Preload("ProductImages").Preload("TrashDetail").Where("id = ?", productID).First(&product).Error
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}
