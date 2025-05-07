package services

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
)

type RequestPickupService interface {
	CreateRequestPickup(request dto.RequestPickup, UserId string) (*dto.ResponseRequestPickup, error)
	GetRequestPickupByID(id string) (*dto.ResponseRequestPickup, error)
	GetAllRequestPickups() ([]dto.ResponseRequestPickup, error)
	UpdateRequestPickup(id string, request dto.RequestPickup) (*dto.ResponseRequestPickup, error)
	DeleteRequestPickup(id string) error
}

type requestPickupService struct {
	repo        repositories.RequestPickupRepository
	repoAddress repositories.AddressRepository
	repoTrash   repositories.TrashRepository
}

func NewRequestPickupService(repo repositories.RequestPickupRepository,
	repoAddress repositories.AddressRepository,
	repoTrash repositories.TrashRepository) RequestPickupService {
	return &requestPickupService{repo: repo, repoAddress: repoAddress, repoTrash: repoTrash}
}

func (s *requestPickupService) CreateRequestPickup(request dto.RequestPickup, UserId string) (*dto.ResponseRequestPickup, error) {

	errors, valid := request.ValidateRequestPickup()
	if !valid {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	findAddress, err := s.repoAddress.FindAddressByID(request.AddressID)
	if err != nil {
		return nil, fmt.Errorf("address with ID %s not found", request.AddressID)
	}

	existingRequest, err := s.repo.FindRequestPickupByAddressAndStatus(UserId, "waiting_pengepul")
	if err != nil {
		return nil, fmt.Errorf("error checking for existing request pickup: %v", err)
	}
	if existingRequest != nil {
		return nil, fmt.Errorf("there is already a pending pickup request for address %s", request.AddressID)
	}

	modelRequest := model.RequestPickup{
		UserId:        UserId,
		AddressId:     findAddress.ID,
		EvidenceImage: request.EvidenceImage,
	}

	err = s.repo.CreateRequestPickup(&modelRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create request pickup: %v", err)
	}

	response := &dto.ResponseRequestPickup{
		ID:            modelRequest.ID,
		UserId:        UserId,
		AddressID:     modelRequest.AddressId,
		EvidenceImage: modelRequest.EvidenceImage,
		StatusPickup:  modelRequest.StatusPickup,
		CreatedAt:     modelRequest.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     modelRequest.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	for _, item := range request.RequestItems {

		findTrashCategory, err := s.repoTrash.GetCategoryByID(item.TrashCategoryID)
		if err != nil {
			return nil, fmt.Errorf("trash category with ID %s not found", item.TrashCategoryID)
		}

		modelItem := model.RequestPickupItem{
			RequestPickupId: modelRequest.ID,
			TrashCategoryId: findTrashCategory.ID,
			EstimatedAmount: item.EstimatedAmount,
		}
		err = s.repo.CreateRequestPickupItem(&modelItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create request pickup item: %v", err)
		}

		response.RequestItems = append(response.RequestItems, dto.ResponseRequestPickupItem{
			ID:                modelItem.ID,
			TrashCategoryName: findTrashCategory.Name,
			EstimatedAmount:   modelItem.EstimatedAmount,
		})
	}

	return response, nil
}

func (s *requestPickupService) GetRequestPickupByID(id string) (*dto.ResponseRequestPickup, error) {

	request, err := s.repo.FindRequestPickupByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching request pickup with ID %s: %v", id, err)
	}

	response := &dto.ResponseRequestPickup{
		ID:            request.ID,
		UserId:        request.UserId,
		AddressID:     request.AddressId,
		EvidenceImage: request.EvidenceImage,
		StatusPickup:  request.StatusPickup,
		CreatedAt:     request.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     request.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *requestPickupService) GetAllRequestPickups() ([]dto.ResponseRequestPickup, error) {

	requests, err := s.repo.FindAllRequestPickups()
	if err != nil {
		return nil, fmt.Errorf("error fetching all request pickups: %v", err)
	}

	var response []dto.ResponseRequestPickup
	for _, request := range requests {
		response = append(response, dto.ResponseRequestPickup{
			ID:            request.ID,
			UserId:        request.UserId,
			AddressID:     request.AddressId,
			EvidenceImage: request.EvidenceImage,
			StatusPickup:  request.StatusPickup,
			CreatedAt:     request.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     request.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}

func (s *requestPickupService) UpdateRequestPickup(id string, request dto.RequestPickup) (*dto.ResponseRequestPickup, error) {

	errors, valid := request.ValidateRequestPickup()
	if !valid {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	existingRequest, err := s.repo.FindRequestPickupByID(id)
	if err != nil {
		return nil, fmt.Errorf("request pickup with ID %s not found: %v", id, err)
	}

	existingRequest.EvidenceImage = request.EvidenceImage

	err = s.repo.UpdateRequestPickup(id, existingRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to update request pickup: %v", err)
	}

	response := &dto.ResponseRequestPickup{
		ID:            existingRequest.ID,
		UserId:        existingRequest.UserId,
		AddressID:     existingRequest.AddressId,
		EvidenceImage: existingRequest.EvidenceImage,
		StatusPickup:  existingRequest.StatusPickup,
		CreatedAt:     existingRequest.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     existingRequest.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *requestPickupService) DeleteRequestPickup(id string) error {

	err := s.repo.DeleteRequestPickup(id)
	if err != nil {
		return fmt.Errorf("failed to delete request pickup with ID %s: %v", id, err)
	}

	return nil
}
