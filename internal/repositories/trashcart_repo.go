package repositories

import (
	"rijig/config"
	"rijig/model"
)

type CartRepository interface {
	CreateCart(cart *model.Cart) error
	GetTrashCategoryByID(id string) (*model.TrashCategory, error)
	GetCartByUserID(userID string) (*model.Cart, error)
	DeleteCartByUserID(userID string) error
}

type cartRepository struct{}

func NewCartRepository() CartRepository {
	return &cartRepository{}
}

func (r *cartRepository) CreateCart(cart *model.Cart) error {
	return config.DB.Create(cart).Error
}

func (r *cartRepository) DeleteCartByUserID(userID string) error {
	return config.DB.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

func (r *cartRepository) GetTrashCategoryByID(id string) (*model.TrashCategory, error) {
	var trash model.TrashCategory
	if err := config.DB.First(&trash, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &trash, nil
}

func (r *cartRepository) GetCartByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := config.DB.Preload("CartItems.TrashCategory").
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
