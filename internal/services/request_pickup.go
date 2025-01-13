package services

import (
	"fmt"

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

func (s *RequestPickupService) GetRequestPickupByID(id string) (*domain.RequestPickup, error) {
	return s.repository.GetByID(id)
}

func (s *RequestPickupService) GetRequestPickupsByUser(userID string) ([]domain.RequestPickup, error) {
	return s.repository.GetByUserID(userID)
}

func (s *RequestPickupService) DeleteRequestPickupByID(id string) error {

	exists, err := s.repository.ExistsByID(id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("request pickup with id %s not found", id)
	}

	return s.repository.DeleteByID(id)
}
