package services

import (
	"context"
	"time"
	"rijig/model"
	"rijig/internal/repositories"
)

type PickupStatusHistoryService interface {
	LogStatusChange(ctx context.Context, requestID, status, changedByID, changedByRole string) error
	GetStatusHistory(ctx context.Context, requestID string) ([]model.PickupStatusHistory, error)
}

type pickupStatusHistoryService struct {
	repo repositories.PickupStatusHistoryRepository
}

func NewPickupStatusHistoryService(repo repositories.PickupStatusHistoryRepository) PickupStatusHistoryService {
	return &pickupStatusHistoryService{repo: repo}
}

func (s *pickupStatusHistoryService) LogStatusChange(ctx context.Context, requestID, status, changedByID, changedByRole string) error {
	history := model.PickupStatusHistory{
		RequestID:     requestID,
		Status:        status,
		ChangedAt:     time.Now(),
		ChangedByID:   changedByID,
		ChangedByRole: changedByRole,
	}
	return s.repo.CreateStatusHistory(ctx, history)
}

func (s *pickupStatusHistoryService) GetStatusHistory(ctx context.Context, requestID string) ([]model.PickupStatusHistory, error) {
	return s.repo.GetStatusHistoryByRequestID(ctx, requestID)
}
