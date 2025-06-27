package requestpickup

import (
	"context"
	"rijig/config"
	"rijig/model"
)

type PickupStatusHistoryRepository interface {
	CreateStatusHistory(ctx context.Context, history model.PickupStatusHistory) error
	GetStatusHistoryByRequestID(ctx context.Context, requestID string) ([]model.PickupStatusHistory, error)
}

type pickupStatusHistoryRepository struct{}

func NewPickupStatusHistoryRepository() PickupStatusHistoryRepository {
	return &pickupStatusHistoryRepository{}
}

func (r *pickupStatusHistoryRepository) CreateStatusHistory(ctx context.Context, history model.PickupStatusHistory) error {
	return config.DB.WithContext(ctx).Create(&history).Error
}

func (r *pickupStatusHistoryRepository) GetStatusHistoryByRequestID(ctx context.Context, requestID string) ([]model.PickupStatusHistory, error) {
	var histories []model.PickupStatusHistory
	err := config.DB.WithContext(ctx).
		Where("request_id = ?", requestID).
		Order("changed_at asc").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}
