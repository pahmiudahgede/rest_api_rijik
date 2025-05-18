package repositories

import (
	"errors"
	"log"

	"rijig/config"
	"rijig/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *model.Cart) error
	GetByUserID(userID string) (*model.Cart, error)
	Update(cart *model.Cart) error
	InsertCartItem(item *model.CartItem) error
	UpdateCartItem(item *model.CartItem) error
	DeleteCartItemByID(id string) error
	Delete(cartID string) error
	DeleteByUserID(userID string) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository() CartRepository {
	return &cartRepository{
		db: config.DB,
	}
}

func (r *cartRepository) Create(cart *model.Cart) error {
	tx := r.db.Begin()
	if err := tx.Create(cart).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *cartRepository) GetByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart

	err := r.db.
		Preload("CartItems.TrashCategory").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Error retrieving cart for user %s: %v", userID, err)
		return nil, errors.New("failed to retrieve cart")
	}

	return &cart, nil
}

func (r *cartRepository) Update(cart *model.Cart) error {
	err := r.db.Save(cart).Error
	if err != nil {
		log.Printf("Error updating cart %s: %v", cart.ID, err)
		return errors.New("failed to update cart")
	}
	return nil
}

func (r *cartRepository) InsertCartItem(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateCartItem(item *model.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) DeleteCartItemByID(id string) error {
	return r.db.Delete(&model.CartItem{}, "id = ?", id).Error
}

func (r *cartRepository) Delete(cartID string) error {
	result := r.db.Where("id = ?", cartID).Delete(&model.Cart{})
	if result.Error != nil {
		log.Printf("Error deleting cart %s: %v", cartID, result.Error)
		return errors.New("failed to delete cart")
	}

	if result.RowsAffected == 0 {
		log.Printf("Cart with ID %s not found for deletion", cartID)
		return errors.New("cart not found")
	}

	return nil
}

func (r *cartRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}
