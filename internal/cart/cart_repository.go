package cart

import (
	"context"
	"errors"
	"fmt"

	"rijig/config"
	"rijig/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindOrCreateCart(ctx context.Context, userID string) (*model.Cart, error)
	AddOrUpdateCartItem(ctx context.Context, cartID, trashCategoryID string, amount float64, estimatedPrice float64) error
	DeleteCartItem(ctx context.Context, cartID, trashCategoryID string) error
	GetCartByUser(ctx context.Context, userID string) (*model.Cart, error)
	UpdateCartTotals(ctx context.Context, cartID string) error
	DeleteCart(ctx context.Context, userID string) error

	CreateCartWithItems(ctx context.Context, cart *model.Cart) error
	HasExistingCart(ctx context.Context, userID string) (bool, error)
}

type cartRepository struct{}

func NewCartRepository() CartRepository {
	return &cartRepository{}
}

func (r *cartRepository) FindOrCreateCart(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	db := config.DB.WithContext(ctx)

	err := db.
		Preload("CartItems.TrashCategory").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err == nil {
		return &cart, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newCart := model.Cart{
			UserID:              userID,
			TotalAmount:         0,
			EstimatedTotalPrice: 0,
		}
		if err := db.Create(&newCart).Error; err != nil {
			return nil, err
		}
		return &newCart, nil
	}

	return nil, err
}

func (r *cartRepository) AddOrUpdateCartItem(ctx context.Context, cartID, trashCategoryID string, amount float64, estimatedPrice float64) error {
	db := config.DB.WithContext(ctx)

	var item model.CartItem
	err := db.
		Where("cart_id = ? AND trash_category_id = ?", cartID, trashCategoryID).
		First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := model.CartItem{
			CartID:                 cartID,
			TrashCategoryID:        trashCategoryID,
			Amount:                 amount,
			SubTotalEstimatedPrice: amount * estimatedPrice,
		}
		return db.Create(&newItem).Error
	}

	if err != nil {
		return err
	}

	item.Amount = amount
	item.SubTotalEstimatedPrice = amount * estimatedPrice
	return db.Save(&item).Error
}

func (r *cartRepository) DeleteCartItem(ctx context.Context, cartID, trashCategoryID string) error {
	db := config.DB.WithContext(ctx)
	return db.Where("cart_id = ? AND trash_category_id = ?", cartID, trashCategoryID).
		Delete(&model.CartItem{}).Error
}

func (r *cartRepository) GetCartByUser(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	db := config.DB.WithContext(ctx)

	err := db.
		Preload("CartItems.TrashCategory").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) UpdateCartTotals(ctx context.Context, cartID string) error {
	db := config.DB.WithContext(ctx)

	var items []model.CartItem
	if err := db.Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
		return err
	}

	var totalAmount float64
	var totalPrice float64

	for _, item := range items {
		totalAmount += item.Amount
		totalPrice += item.SubTotalEstimatedPrice
	}

	return db.Model(&model.Cart{}).
		Where("id = ?", cartID).
		Updates(map[string]interface{}{
			"total_amount":          totalAmount,
			"estimated_total_price": totalPrice,
		}).Error
}

func (r *cartRepository) DeleteCart(ctx context.Context, userID string) error {
	db := config.DB.WithContext(ctx)
	var cart model.Cart
	if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}
	return db.Delete(&cart).Error
}

func (r *cartRepository) CreateCartWithItems(ctx context.Context, cart *model.Cart) error {
	db := config.DB.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cart).Error; err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
		return nil
	})
}

func (r *cartRepository) HasExistingCart(ctx context.Context, userID string) (bool, error) {
	db := config.DB.WithContext(ctx)

	var count int64
	err := db.Model(&model.Cart{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
