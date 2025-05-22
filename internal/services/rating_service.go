package services

import (
	"context"
	"time"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
)

type PickupRatingService interface {
	CreateRating(ctx context.Context, userID, pickupID, collectorID string, input dto.CreatePickupRatingDTO) error
	GetRatingsByCollector(ctx context.Context, collectorID string) ([]model.PickupRating, error)
	GetAverageRating(ctx context.Context, collectorID string) (float32, error)
}

type pickupRatingService struct {
	repo repositories.PickupRatingRepository
}

func NewPickupRatingService(repo repositories.PickupRatingRepository) PickupRatingService {
	return &pickupRatingService{repo: repo}
}

func (s *pickupRatingService) CreateRating(ctx context.Context, userID, pickupID, collectorID string, input dto.CreatePickupRatingDTO) error {
	rating := model.PickupRating{
		RequestID:   pickupID,
		UserID:      userID,
		CollectorID: collectorID,
		Rating:      input.Rating,
		Feedback:    input.Feedback,
		CreatedAt:   time.Now(),
	}
	return s.repo.CreateRating(ctx, rating)
}

func (s *pickupRatingService) GetRatingsByCollector(ctx context.Context, collectorID string) ([]model.PickupRating, error) {
	return s.repo.GetRatingsByCollector(ctx, collectorID)
}

func (s *pickupRatingService) GetAverageRating(ctx context.Context, collectorID string) (float32, error) {
	return s.repo.CalculateAverageRating(ctx, collectorID)
}
