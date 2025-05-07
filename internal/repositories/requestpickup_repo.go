package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type RequestPickupRepository interface {
	CreateRequestPickup(request *model.RequestPickup) error
	CreateRequestPickupItem(item *model.RequestPickupItem) error
	FindRequestPickupByID(id string) (*model.RequestPickup, error)
	FindAllRequestPickups() ([]model.RequestPickup, error)
	UpdateRequestPickup(id string, request *model.RequestPickup) error
	DeleteRequestPickup(id string) error
	FindRequestPickupByAddressAndCategory(addressID string, trashCategoryID string) (*model.RequestPickup, error)
	FindRequestPickupByAddressAndStatus(userId, status string) (*model.RequestPickup, error)
}

type requestPickupRepository struct {
	DB *gorm.DB
}

func NewRequestPickupRepository(db *gorm.DB) RequestPickupRepository {
	return &requestPickupRepository{DB: db}
}

func (r *requestPickupRepository) CreateRequestPickup(request *model.RequestPickup) error {
	if err := r.DB.Create(request).Error; err != nil {
		return fmt.Errorf("failed to create request pickup: %v", err)
	}

	for _, item := range request.RequestItems {
		item.RequestPickupId = request.ID
		if err := r.DB.Create(&item).Error; err != nil {
			return fmt.Errorf("failed to create request pickup item: %v", err)
		}
	}

	return nil
}

func (r *requestPickupRepository) CreateRequestPickupItem(item *model.RequestPickupItem) error {
	if err := r.DB.Create(item).Error; err != nil {
		return fmt.Errorf("failed to create request pickup item: %v", err)
	}
	return nil
}

func (r *requestPickupRepository) FindRequestPickupByID(id string) (*model.RequestPickup, error) {
	var request model.RequestPickup
	err := r.DB.Preload("RequestItems").First(&request, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("request pickup with ID %s not found: %v", id, err)
	}
	return &request, nil
}

func (r *requestPickupRepository) FindAllRequestPickups() ([]model.RequestPickup, error) {
	var requests []model.RequestPickup
	err := r.DB.Preload("RequestItems").Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all request pickups: %v", err)
	}
	return requests, nil
}

func (r *requestPickupRepository) UpdateRequestPickup(id string, request *model.RequestPickup) error {
	err := r.DB.Model(&model.RequestPickup{}).Where("id = ?", id).Updates(request).Error
	if err != nil {
		return fmt.Errorf("failed to update request pickup: %v", err)
	}
	return nil
}

func (r *requestPickupRepository) DeleteRequestPickup(id string) error {

	if err := r.DB.Where("request_pickup_id = ?", id).Delete(&model.RequestPickupItem{}).Error; err != nil {
		return fmt.Errorf("failed to delete request pickup items: %v", err)
	}

	err := r.DB.Delete(&model.RequestPickup{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("failed to delete request pickup: %v", err)
	}
	return nil
}

func (r *requestPickupRepository) FindRequestPickupByAddressAndCategory(addressID string, trashCategoryID string) (*model.RequestPickup, error) {
	var request model.RequestPickup
	err := r.DB.Joins("JOIN request_pickup_items ON request_pickups.id = request_pickup_items.request_pickup_id").
		Where("request_pickups.address_id = ? AND request_pickup_items.trash_category_id = ?", addressID, trashCategoryID).
		First(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error checking request pickup for address %s and category %s: %v", addressID, trashCategoryID, err)
	}
	return &request, nil
}

func (r *requestPickupRepository) FindRequestPickupByAddressAndStatus(userId, status string) (*model.RequestPickup, error) {
	var request model.RequestPickup
	err := r.DB.Where("user_id = ? AND status_pickup = ?", userId, status).First(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check existing request pickup: %v", err)
	}
	return &request, nil
}
