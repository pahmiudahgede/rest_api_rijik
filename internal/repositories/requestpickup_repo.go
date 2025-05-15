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
	FindAllRequestPickups(userId string) ([]model.RequestPickup, error)
	FindAllAutomaticMethodRequest(requestMethod, statuspickup string) ([]model.RequestPickup, error)
	FindRequestPickupByAddressAndStatus(userId, status, method string) (*model.RequestPickup, error)
	FindRequestPickupByStatus(userId, status, method string) (*model.RequestPickup, error)
	GetRequestPickupItems(requestPickupId string) ([]model.RequestPickupItem, error)
	GetAutomaticRequestPickupsForCollector() ([]model.RequestPickup, error)
	GetManualReqMethodforCollect(collector_id string) ([]model.RequestPickup, error)
	// SelectCollectorInRequest(userId string, collectorId string) error
	UpdateRequestPickup(id string, request *model.RequestPickup) error
	DeleteRequestPickup(id string) error
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

func (r *requestPickupRepository) FindAllRequestPickups(userId string) ([]model.RequestPickup, error) {
	var requests []model.RequestPickup
	err := r.DB.Preload("RequestItems").Where("user_id = ?", userId).Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all request pickups: %v", err)
	}
	return requests, nil
}

func (r *requestPickupRepository) FindAllAutomaticMethodRequest(requestMethod, statuspickup string) ([]model.RequestPickup, error) {
	var requests []model.RequestPickup
	err := r.DB.Preload("RequestItems").Where("request_method = ? AND status_pickup = ?", requestMethod, statuspickup).Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching request pickups with request_method %s: %v", requestMethod, err)
	}

	return requests, nil
}

func (r *requestPickupRepository) FindRequestPickupByAddressAndStatus(userId, status, method string) (*model.RequestPickup, error) {
	var request model.RequestPickup
	err := r.DB.Preload("Address").Where("user_id = ? AND status_pickup = ? AND request_method =?", userId, status, method).First(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check existing request pickup: %v", err)
	}
	return &request, nil
}

func (r *requestPickupRepository) FindRequestPickupByStatus(userId, status, method string) (*model.RequestPickup, error) {
	var request model.RequestPickup
	err := r.DB.Where("user_id = ? AND status_pickup = ? AND request_method =?", userId, status, method).First(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check existing request pickup: %v", err)
	}
	return &request, nil
}

func (r *requestPickupRepository) UpdateRequestPickup(id string, request *model.RequestPickup) error {
	err := r.DB.Model(&model.RequestPickup{}).Where("id = ?", id).Updates(request).Error
	if err != nil {
		return fmt.Errorf("failed to update request pickup: %v", err)
	}

	return nil
}

// func (r *requestPickupRepository) SelectCollectorInRequest(userId string, collectorId string) error {
// 	var request model.RequestPickup
// 	err := r.DB.Model(&model.RequestPickup{}).
// 		Where("user_id = ? AND status_pickup = ? AND request_method = ? AND collector_id IS NULL", userId, "waiting_collector", "manual").
// 		First(&request).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return fmt.Errorf("no matching request pickup found for user %s", userId)
// 		}
// 		return fmt.Errorf("failed to find request pickup: %v", err)
// 	}

// 	err = r.DB.Model(&model.RequestPickup{}).
// 		Where("id = ?", request.ID).
// 		Update("collector_id", collectorId).
// 		Error
// 	if err != nil {
// 		return fmt.Errorf("failed to update collector_id: %v", err)
// 	}
// 	return nil
// }

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

func (r *requestPickupRepository) GetAutomaticRequestPickupsForCollector() ([]model.RequestPickup, error) {
	var requests []model.RequestPickup
	err := r.DB.Preload("Address").
		Where("request_method = ? AND status_pickup = ? AND collector_id = ?", "otomatis", "waiting_collector", nil).
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching pickup requests: %v", err)
	}
	return requests, nil
}

func (r *requestPickupRepository) GetManualReqMethodforCollect(collector_id string) ([]model.RequestPickup, error) {
	var requests []model.RequestPickup
	err := r.DB.Where("request_method = ? AND status_pickup = ? AND collector_id = ?", "otomatis", "waiting_collector", collector_id).
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching pickup requests: %v", err)
	}
	return requests, nil
}

func (r *requestPickupRepository) GetRequestPickupItems(requestPickupId string) ([]model.RequestPickupItem, error) {
	var items []model.RequestPickupItem

	err := r.DB.Preload("TrashCategory").Where("request_pickup_id = ?", requestPickupId).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching request pickup items: %v", err)
	}
	return items, nil
}
