package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

type RequestPickupRepository interface {
	Create(request *domain.RequestPickup) error
	GetByID(id string) (*domain.RequestPickup, error)
	GetByUserID(userID string) ([]domain.RequestPickup, error)
}

type requestPickupRepository struct{}

func NewRequestPickupRepository() RequestPickupRepository {
	return &requestPickupRepository{}
}

func (r *requestPickupRepository) Create(request *domain.RequestPickup) error {
	return config.DB.Create(request).Error
}

func (r *requestPickupRepository) GetByID(id string) (*domain.RequestPickup, error) {
	var requestPickup domain.RequestPickup
	if err := config.DB.Preload("Request").
		Preload("Request.TrashCategory").
		Preload("UserAddress").
		Where("id = ?", id).
		First(&requestPickup).Error; err != nil {
		return nil, err
	}
	return &requestPickup, nil
}

func (r *requestPickupRepository) GetByUserID(userID string) ([]domain.RequestPickup, error) {
	var requestPickups []domain.RequestPickup
	err := config.DB.Preload("Request").
		Preload("Request.TrashCategory").
		Preload("UserAddress").
		Where("user_id = ?", userID).
		Find(&requestPickups).Error

	if err != nil {
		return nil, err
	}

	return requestPickups, nil
}
