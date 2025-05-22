package repositories

import (
	"context"
	"rijig/config"
	"rijig/dto"
	"rijig/model"
	"time"
)

type RequestPickupRepository interface {
	CreateRequestPickup(ctx context.Context, pickup *model.RequestPickup) error
	GetPickupWithItemsAndAddress(ctx context.Context, id string) (*model.RequestPickup, error)
	GetAllAutomaticRequestsWithAddress(ctx context.Context) ([]model.RequestPickup, error)
	UpdateCollectorID(ctx context.Context, pickupID, collectorID string) error
	GetRequestsAssignedToCollector(ctx context.Context, collectorID string) ([]model.RequestPickup, error)
	UpdatePickupStatusAndConfirmationTime(ctx context.Context, pickupID string, status string, confirmedAt time.Time) error
	UpdatePickupStatus(ctx context.Context, pickupID string, status string) error
	UpdateRequestPickupItemsAmountAndPrice(ctx context.Context, pickupID string, items []dto.UpdateRequestPickupItemDTO) error
}

type requestPickupRepository struct{}

func NewRequestPickupRepository() RequestPickupRepository {
	return &requestPickupRepository{}
}

func (r *requestPickupRepository) CreateRequestPickup(ctx context.Context, pickup *model.RequestPickup) error {
	return config.DB.WithContext(ctx).Create(pickup).Error
}

func (r *requestPickupRepository) GetPickupWithItemsAndAddress(ctx context.Context, id string) (*model.RequestPickup, error) {
	var pickup model.RequestPickup
	err := config.DB.WithContext(ctx).
		Preload("RequestItems").
		Preload("Address").
		Where("id = ?", id).
		First(&pickup).Error

	if err != nil {
		return nil, err
	}
	return &pickup, nil
}

func (r *requestPickupRepository) UpdateCollectorID(ctx context.Context, pickupID, collectorID string) error {
	return config.DB.WithContext(ctx).
		Model(&model.RequestPickup{}).
		Where("id = ?", pickupID).
		Update("collector_id", collectorID).
		Error
}

func (r *requestPickupRepository) GetAllAutomaticRequestsWithAddress(ctx context.Context) ([]model.RequestPickup, error) {
	var pickups []model.RequestPickup
	err := config.DB.WithContext(ctx).
		Preload("RequestItems").
		Preload("Address").
		Where("request_method = ?", "otomatis").
		Find(&pickups).Error

	if err != nil {
		return nil, err
	}
	return pickups, nil
}

func (r *requestPickupRepository) GetRequestsAssignedToCollector(ctx context.Context, collectorID string) ([]model.RequestPickup, error) {
	var pickups []model.RequestPickup
	err := config.DB.WithContext(ctx).
		Preload("User").
		Preload("Address").
		Preload("RequestItems").
		Where("collector_id = ? AND status_pickup = ?", collectorID, "waiting_collector").
		Find(&pickups).Error

	if err != nil {
		return nil, err
	}
	return pickups, nil
}

func (r *requestPickupRepository) UpdatePickupStatusAndConfirmationTime(ctx context.Context, pickupID string, status string, confirmedAt time.Time) error {
	return config.DB.WithContext(ctx).
		Model(&model.RequestPickup{}).
		Where("id = ?", pickupID).
		Updates(map[string]interface{}{
			"status_pickup":             status,
			"confirmed_by_collector_at": confirmedAt,
		}).Error
}

func (r *requestPickupRepository) UpdatePickupStatus(ctx context.Context, pickupID string, status string) error {
	return config.DB.WithContext(ctx).
		Model(&model.RequestPickup{}).
		Where("id = ?", pickupID).
		Update("status_pickup", status).
		Error
}

func (r *requestPickupRepository) UpdateRequestPickupItemsAmountAndPrice(ctx context.Context, pickupID string, items []dto.UpdateRequestPickupItemDTO) error {
	// ambil collector_id dulu dari pickup
	var pickup model.RequestPickup
	if err := config.DB.WithContext(ctx).
		Select("collector_id").
		Where("id = ?", pickupID).
		First(&pickup).Error; err != nil {
		return err
	}

	for _, item := range items {
		var pickupItem model.RequestPickupItem
		err := config.DB.WithContext(ctx).
			Where("id = ? AND request_pickup_id = ?", item.ItemID, pickupID).
			First(&pickupItem).Error
		if err != nil {
			return err
		}

		var price float64
		err = config.DB.WithContext(ctx).
			Model(&model.AvaibleTrashByCollector{}).
			Where("collector_id = ? AND trash_category_id = ?", pickup.CollectorID, pickupItem.TrashCategoryId).
			Select("price").
			Scan(&price).Error
		if err != nil {
			return err
		}

		finalPrice := item.Amount * price
		err = config.DB.WithContext(ctx).
			Model(&model.RequestPickupItem{}).
			Where("id = ?", item.ItemID).
			Updates(map[string]interface{}{
				"estimated_amount": item.Amount,
				"final_price":      finalPrice,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
