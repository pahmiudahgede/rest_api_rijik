package repositories

import (
	"fmt"

	"rijig/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CountProductsByStoreID(storeID string) (int64, error)
	CreateProduct(product *model.Product) error
	GetProductByID(productID string) (*model.Product, error)
	GetProductsByStoreID(storeID string) ([]model.Product, error)
	FindProductsByStoreID(storeID string, page, limit int) ([]model.Product, error)
	FindProductImagesByProductID(productID string) ([]model.ProductImage, error)
	GetProductImageByID(imageID string) (*model.ProductImage, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(productID string) error
	DeleteProductsByID(productIDs []string) error
	AddProductImages(images []model.ProductImage) error
	DeleteProductImagesByProductID(productID string) error
	DeleteProductImagesByID(imageIDs []string) error
	DeleteProductImageByID(imageID string) error
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{DB}
}

func (r *productRepository) CreateProduct(product *model.Product) error {
	return r.DB.Create(product).Error
}

func (r *productRepository) CountProductsByStoreID(storeID string) (int64, error) {
	var count int64
	if err := r.DB.Model(&model.Product{}).Where("store_id = ?", storeID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) GetProductByID(productID string) (*model.Product, error) {
	var product model.Product
	if err := r.DB.Preload("ProductImages").Where("id = ?", productID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetProductsByStoreID(storeID string) ([]model.Product, error) {
	var products []model.Product
	if err := r.DB.Where("store_id = ?", storeID).Preload("ProductImages").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindProductsByStoreID(storeID string, page, limit int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * limit

	if err := r.DB.
		Where("store_id = ?", storeID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FindProductImagesByProductID(productID string) ([]model.ProductImage, error) {
	var productImages []model.ProductImage
	if err := r.DB.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return nil, err
	}
	return productImages, nil
}

func (r *productRepository) GetProductImageByID(imageID string) (*model.ProductImage, error) {
	var productImage model.ProductImage
	if err := r.DB.Where("id = ?", imageID).First(&productImage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &productImage, nil
}

func (r *productRepository) UpdateProduct(product *model.Product) error {
	return r.DB.Save(product).Error
}

func (r *productRepository) DeleteProduct(productID string) error {
	return r.DB.Delete(&model.Product{}, "id = ?", productID).Error
}

func (r *productRepository) DeleteProductsByID(productIDs []string) error {
	if err := r.DB.Where("id IN ?", productIDs).Delete(&model.Product{}).Error; err != nil {
		return fmt.Errorf("failed to delete products: %v", err)
	}
	return nil
}

func (r *productRepository) AddProductImages(images []model.ProductImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.DB.Create(&images).Error
}

func (r *productRepository) DeleteProductImagesByProductID(productID string) error {
	return r.DB.Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error
}

func (r *productRepository) DeleteProductImagesByID(imageIDs []string) error {
	if err := r.DB.Where("id IN ?", imageIDs).Delete(&model.ProductImage{}).Error; err != nil {
		return fmt.Errorf("failed to delete product images: %v", err)
	}
	return nil
}

func (r *productRepository) DeleteProductImageByID(imageID string) error {
	if err := r.DB.Where("id = ?", imageID).Delete(&model.ProductImage{}).Error; err != nil {
		return fmt.Errorf("failed to delete product image: %v", err)
	}
	return nil
}
