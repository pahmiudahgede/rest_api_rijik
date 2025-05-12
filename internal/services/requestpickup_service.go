package services

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type RequestPickupService interface {
	CreateRequestPickup(request dto.RequestPickup, UserId string) (*dto.ResponseRequestPickup, error)
	GetRequestPickupByID(id string) (*dto.ResponseRequestPickup, error)
	GetAllRequestPickups(userid string) ([]dto.ResponseRequestPickup, error)
	// GetAllAutomaticRequestPickups(collector_id string) ([]dto.ResponseRequestPickup, error)
	// GetAllAutomaticRequestPickup(collectorId string) ([]dto.ResponseRequestPickup, error)

	GetRequestPickupsForCollector(collectorId string) ([]dto.ResponseRequestPickup, error)
}

type requestPickupService struct {
	repo        repositories.RequestPickupRepository
	repoReq     repositories.CollectorRepository
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

	_, err := s.repoAddress.FindAddressByID(request.AddressID)
	if err != nil {
		return nil, fmt.Errorf("address with ID %s not found", request.AddressID)
	}

	existingRequest, err := s.repo.FindRequestPickupByAddressAndStatus(UserId, "waiting_collector")
	if err != nil {
		return nil, fmt.Errorf("error checking for existing request pickup: %v", err)
	}
	if existingRequest != nil {
		return nil, fmt.Errorf("there is already a pending pickup request for address %s", request.AddressID)
	}

	modelRequest := model.RequestPickup{
		UserId:        UserId,
		AddressId:     request.AddressID,
		EvidenceImage: request.EvidenceImage,
		RequestMethod: request.RequestMethod,
	}

	err = s.repo.CreateRequestPickup(&modelRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create request pickup: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(modelRequest.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(modelRequest.UpdatedAt)

	response := &dto.ResponseRequestPickup{
		ID:            modelRequest.ID,
		UserId:        UserId,
		AddressID:     modelRequest.AddressId,
		EvidenceImage: modelRequest.EvidenceImage,
		StatusPickup:  modelRequest.StatusPickup,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
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

	createdAt, _ := utils.FormatDateToIndonesianFormat(request.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(request.UpdatedAt)

	response := &dto.ResponseRequestPickup{
		ID:            request.ID,
		UserId:        request.UserId,
		AddressID:     request.AddressId,
		EvidenceImage: request.EvidenceImage,
		StatusPickup:  request.StatusPickup,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return response, nil
}

func (s *requestPickupService) GetAllRequestPickups(userid string) ([]dto.ResponseRequestPickup, error) {

	requests, err := s.repo.FindAllRequestPickups(userid)
	if err != nil {
		return nil, fmt.Errorf("error fetching all request pickups: %v", err)
	}

	var response []dto.ResponseRequestPickup
	for _, request := range requests {
		createdAt, _ := utils.FormatDateToIndonesianFormat(request.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(request.UpdatedAt)
		response = append(response, dto.ResponseRequestPickup{
			ID:            request.ID,
			UserId:        request.UserId,
			AddressID:     request.AddressId,
			EvidenceImage: request.EvidenceImage,
			StatusPickup:  request.StatusPickup,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	return response, nil
}

func (s *requestPickupService) GetRequestPickupsForCollector(collectorId string) ([]dto.ResponseRequestPickup, error) {

	requests, err := s.repo.GetAutomaticRequestPickupsForCollector(collectorId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving automatic pickup requests: %v", err)
	}

	var response []dto.ResponseRequestPickup

	for _, req := range requests {

		_, distance := utils.Distance(
			utils.Coord{Lat: req.Address.Latitude, Lon: req.Address.Longitude},
			utils.Coord{Lat: req.Address.Latitude, Lon: req.Address.Longitude},
		)

		if distance <= 20 {

			mappedRequest := dto.ResponseRequestPickup{
				ID:            req.ID,
				UserId:        req.UserId,
				AddressID:     req.AddressId,
				EvidenceImage: req.EvidenceImage,
				StatusPickup:  req.StatusPickup,
				CreatedAt:     req.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:     req.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			requestItems, err := s.repo.GetRequestPickupItems(req.ID)
			if err != nil {
				return nil, fmt.Errorf("error fetching request items: %v", err)
			}

			var mappedRequestItems []dto.ResponseRequestPickupItem
			for _, item := range requestItems {
				mappedRequestItems = append(mappedRequestItems, dto.ResponseRequestPickupItem{
					ID:                item.ID,
					TrashCategoryName: item.TrashCategory.Name,
					EstimatedAmount:   item.EstimatedAmount,
				})
			}

			mappedRequest.RequestItems = mappedRequestItems

			response = append(response, mappedRequest)
		}
	}

	return response, nil
}
