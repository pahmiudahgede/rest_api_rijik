package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type StoreService interface {
	CreateStore(userID string, storeDTO dto.RequestStoreDTO, storeLogo *multipart.FileHeader, storeBanner *multipart.FileHeader) (*dto.ResponseStoreDTO, error)
	GetStoreByUserID(userID string) (*dto.ResponseStoreDTO, error)
	UpdateStore(storeID string, storeDTO *dto.RequestStoreDTO, storeLogo *multipart.FileHeader, storeBanner *multipart.FileHeader, userID string) (*dto.ResponseStoreDTO, error)
	DeleteStore(storeID string) error
}

type storeService struct {
	storeRepo repositories.StoreRepository
}

func NewStoreService(storeRepo repositories.StoreRepository) StoreService {
	return &storeService{storeRepo}
}

func (s *storeService) CreateStore(userID string, storeDTO dto.RequestStoreDTO, storeLogo, storeBanner *multipart.FileHeader) (*dto.ResponseStoreDTO, error) {

	if errors, valid := storeDTO.ValidateStoreInput(); !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

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

	storeLogoPath, err := s.saveStoreImage(storeLogo, "logo")
	if err != nil {
		return nil, fmt.Errorf("failed to save store logo: %w", err)
	}

	storeBannerPath, err := s.saveStoreImage(storeBanner, "banner")
	if err != nil {
		return nil, fmt.Errorf("failed to save store banner: %w", err)
	}

	store := model.Store{
		UserID:         userID,
		StoreName:      storeDTO.StoreName,
		StoreLogo:      storeLogoPath,
		StoreBanner:    storeBannerPath,
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

func (s *storeService) GetStoreByUserID(userID string) (*dto.ResponseStoreDTO, error) {

	store, err := s.storeRepo.FindStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving store by user ID: %w", err)
	}
	if store == nil {
		return nil, fmt.Errorf("store not found for user ID: %s", userID)
	}

	createdAt, err := utils.FormatDateToIndonesianFormat(store.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format createdAt: %w", err)
	}

	updatedAt, err := utils.FormatDateToIndonesianFormat(store.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format updatedAt: %w", err)
	}

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

func (s *storeService) UpdateStore(storeID string, storeDTO *dto.RequestStoreDTO, storeLogo, storeBanner *multipart.FileHeader, userID string) (*dto.ResponseStoreDTO, error) {
	store, err := s.storeRepo.FindStoreByID(storeID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving store by ID: %w", err)
	}
	if store == nil {
		return nil, fmt.Errorf("store not found")
	}

	if storeDTO.StoreAddressID == "" {
		return nil, fmt.Errorf("store address ID cannot be empty")
	}

	address, err := s.storeRepo.FindAddressByID(storeDTO.StoreAddressID)
	if err != nil {
		return nil, fmt.Errorf("error validating store address ID: %w", err)
	}
	if address == nil {
		return nil, fmt.Errorf("store address ID not found")
	}

	if storeLogo != nil {
		if err := s.deleteStoreImage(store.StoreLogo); err != nil {
			return nil, fmt.Errorf("failed to delete old store logo: %w", err)
		}
		storeLogoPath, err := s.saveStoreImage(storeLogo, "logo")
		if err != nil {
			return nil, fmt.Errorf("failed to save store logo: %w", err)
		}
		store.StoreLogo = storeLogoPath
	}

	if storeBanner != nil {
		if err := s.deleteStoreImage(store.StoreBanner); err != nil {
			return nil, fmt.Errorf("failed to delete old store banner: %w", err)
		}
		storeBannerPath, err := s.saveStoreImage(storeBanner, "banner")
		if err != nil {
			return nil, fmt.Errorf("failed to save store banner: %w", err)
		}
		store.StoreBanner = storeBannerPath
	}

	store.StoreName = storeDTO.StoreName
	store.StoreInfo = storeDTO.StoreInfo
	store.StoreAddressID = storeDTO.StoreAddressID

	if err := s.storeRepo.UpdateStore(store); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	createdAt, err := utils.FormatDateToIndonesianFormat(store.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format createdAt: %w", err)
	}
	updatedAt, err := utils.FormatDateToIndonesianFormat(store.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to format updatedAt: %w", err)
	}

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

func (s *storeService) DeleteStore(storeID string) error {
	store, err := s.storeRepo.FindStoreByID(storeID)
	if err != nil {
		return fmt.Errorf("error retrieving store by ID: %w", err)
	}
	if store == nil {
		return fmt.Errorf("store not found")
	}

	if store.StoreLogo != "" {
		if err := s.deleteStoreImage(store.StoreLogo); err != nil {
			return fmt.Errorf("failed to delete store logo: %w", err)
		}
	}

	if store.StoreBanner != "" {
		if err := s.deleteStoreImage(store.StoreBanner); err != nil {
			return fmt.Errorf("failed to delete store banner: %w", err)
		}
	}

	if err := s.storeRepo.DeleteStore(storeID); err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	return nil
}

func (s *storeService) saveStoreImage(file *multipart.FileHeader, imageType string) (string, error) {

	imageDir := fmt.Sprintf("./public%s/uploads/store/%s",os.Getenv("BASE_URL"), imageType)
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {

		if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for %s image: %v", imageType, err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(file.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed for %s", imageType)
	}

	fileName := fmt.Sprintf("%s_%s%s", imageType, uuid.New().String(), extension)
	filePath := filepath.Join(imageDir, fileName)

	fileData, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fileData.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create %s image file: %v", imageType, err)
	}
	defer outFile.Close()

	if _, err := outFile.ReadFrom(fileData); err != nil {
		return "", fmt.Errorf("failed to save %s image: %v", imageType, err)
	}

	return filepath.Join("/uploads/store", imageType, fileName), nil
}

func (s *storeService) deleteStoreImage(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	filePath := filepath.Join("./public", imagePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file at %s: %w", filePath, err)
	}

	return nil
}
