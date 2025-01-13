package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

type RequestPickupRepository interface {
	Create(request *domain.RequestPickup) error
	GetByID(id string) (*domain.RequestPickup, error)
	GetByUserID(userID string) ([]domain.RequestPickup, error)
	DeleteByID(id string) error
	ExistsByID(id string) (bool, error) 
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

func (r *requestPickupRepository) ExistsByID(id string) (bool, error) {
	var count int64
	if err := config.DB.Model(&domain.RequestPickup{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}


func (r *requestPickupRepository) DeleteByID(id string) error {
	return config.DB.Where("id = ?", id).Delete(&domain.RequestPickup{}).Error
}
