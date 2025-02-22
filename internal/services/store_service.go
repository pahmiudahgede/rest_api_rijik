package services

import (
	"fmt"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type StoreService interface {
	CreateStore(userID string, storeDTO *dto.RequestStoreDTO) (*dto.ResponseStoreDTO, error)
}

type storeService struct {
	storeRepo repositories.StoreRepository
}

func NewStoreService(storeRepo repositories.StoreRepository) StoreService {
	return &storeService{storeRepo}
}

func (s *storeService) CreateStore(userID string, storeDTO *dto.RequestStoreDTO) (*dto.ResponseStoreDTO, error) {

	existingStore, err := s.storeRepo.FindStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking if user already has a store: %w", err)
	}
	if existingStore != nil {
		return nil, fmt.Errorf("user already has a store")
	}

	address, err := s.storeRepo.FindAddressByID(storeDTO.StoreAddressID)
	if err != nil {
		return nil, fmt.Errorf("error validating store address ID: %w", err)
	}
	if address == nil {
		return nil, fmt.Errorf("store address ID not found")
	}

	store := model.Store{
		UserID:         userID,
		StoreName:      storeDTO.StoreName,
		StoreLogo:      storeDTO.StoreLogo,
		StoreBanner:    storeDTO.StoreBanner,
		StoreInfo:      storeDTO.StoreInfo,
		StoreAddressID: storeDTO.StoreAddressID,
	}

	if err := s.storeRepo.CreateStore(&store); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(store.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(store.UpdatedAt)

	storeResponseDTO := &dto.ResponseStoreDTO{
		ID:             store.ID,
		UserID:         store.UserID,
		StoreName:      store.StoreName,
		StoreLogo:      store.StoreLogo,
		StoreBanner:    store.StoreBanner,
		StoreInfo:      store.StoreInfo,
		StoreAddressID: store.StoreAddressID,
		TotalProduct:   store.TotalProduct,
		Followers:      store.Followers,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	return storeResponseDTO, nil
}
