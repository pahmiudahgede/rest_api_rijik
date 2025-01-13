package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

type RequestPickupService struct {
	repository repositories.RequestPickupRepository
}

func NewRequestPickupService(repository repositories.RequestPickupRepository) *RequestPickupService {
	return &RequestPickupService{repository: repository}
}

func (s *RequestPickupService) CreateRequestPickup(request *domain.RequestPickup) error {
	return s.repository.Create(request)
}

func (s *RequestPickupService) GetRequestPickupsByUser(userID string) ([]domain.RequestPickup, error) {
	return s.repository.GetByUserID(userID)
}
