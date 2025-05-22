package repositories

import (
	"context"
	"rijig/config"
	"rijig/model"
)

type PickupRatingRepository interface {
	CreateRating(ctx context.Context, rating model.PickupRating) error
	GetRatingsByCollector(ctx context.Context, collectorID string) ([]model.PickupRating, error)
	CalculateAverageRating(ctx context.Context, collectorID string) (float32, error)
}

type pickupRatingRepository struct{}

func NewPickupRatingRepository() PickupRatingRepository {
	return &pickupRatingRepository{}
}

func (r *pickupRatingRepository) CreateRating(ctx context.Context, rating model.PickupRating) error {
	return config.DB.WithContext(ctx).Create(&rating).Error
}

func (r *pickupRatingRepository) GetRatingsByCollector(ctx context.Context, collectorID string) ([]model.PickupRating, error) {
	var ratings []model.PickupRating
	err := config.DB.WithContext(ctx).
		Where("collector_id = ?", collectorID).
		Order("created_at desc").
		Find(&ratings).Error
	if err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *pickupRatingRepository) CalculateAverageRating(ctx context.Context, collectorID string) (float32, error) {
	var avg float32
	err := config.DB.WithContext(ctx).
		Model(&model.PickupRating{}).
		Select("AVG(rating)").
		Where("collector_id = ?", collectorID).
		Scan(&avg).Error
	if err != nil {
		return 0, err
	}
	return avg, nil
}
