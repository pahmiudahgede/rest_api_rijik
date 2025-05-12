package services

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/utils"
)

type CollectorService interface {
	FindCollectorsNearby(userId string) ([]dto.ResponseCollectorDTO, error)
	ConfirmRequestPickup(requestId, collectorId string) (*dto.ResponseRequestPickup, error)
}

type collectorService struct {
	repo        repositories.CollectorRepository
	repoReq     repositories.RequestPickupRepository
	repoAddress repositories.AddressRepository
}

func NewCollectorService(repo repositories.CollectorRepository,
	repoReq repositories.RequestPickupRepository,
	repoAddress repositories.AddressRepository) CollectorService {
	return &collectorService{repo: repo, repoReq: repoReq, repoAddress: repoAddress}
}

func (s *collectorService) FindCollectorsNearby(userId string) ([]dto.ResponseCollectorDTO, error) {
	collectors, err := s.repo.FindActiveCollectors()
	if err != nil {
		return nil, fmt.Errorf("error fetching active collectors: %v", err)
	}

	request, err := s.repoReq.FindRequestPickupByAddressAndStatus(userId, "waiting_collector")
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data request pickup dengan userid: %v", err)
	}

	reqpickaddress, err := s.repoAddress.FindAddressByID(request.AddressId)
	if err != nil {
		return nil, fmt.Errorf("error fetching address for request pickup %s: %v", request.ID, err)
	}

	var nearbyCollectorsResponse []dto.ResponseCollectorDTO
	var maxDistance = 10.0

	for _, collector := range collectors {

		address, err := s.repoAddress.FindAddressByID(collector.AddressId)
		if err != nil {
			return nil, fmt.Errorf("error fetching address for collector %s: %v", collector.ID, err)
		}

		collectorCoord := utils.Coord{Lat: reqpickaddress.Latitude, Lon: reqpickaddress.Longitude}
		userCoord := utils.Coord{Lat: address.Latitude, Lon: address.Longitude}

		_, km := utils.Distance(collectorCoord, userCoord)

		if km <= maxDistance {

			nearbyCollectorsResponse = append(nearbyCollectorsResponse, dto.ResponseCollectorDTO{
				ID:        collector.ID,
				AddressId: collector.User.Name,
				Rating:    collector.Rating,
			})
		}
	}

	if len(nearbyCollectorsResponse) == 0 {
		return nil, fmt.Errorf("no request pickups found within %v km", maxDistance)
	}

	return nearbyCollectorsResponse, nil
}

func (s *collectorService) ConfirmRequestPickup(requestId, collectorId string) (*dto.ResponseRequestPickup, error) {

	request, err := s.repoReq.FindRequestPickupByID(requestId)
	if err != nil {
		return nil, fmt.Errorf("request pickup not found: %v", err)
	}

	if request.StatusPickup != "waiting_collector" {
		return nil, fmt.Errorf("pickup request is not in 'waiting_collector' status")
	}

	collector, err := s.repo.FindCollectorById(collectorId)
	if err != nil {
		return nil, fmt.Errorf("collector tidak ditemukan: %v", err)
	}

	request.StatusPickup = "confirmed"
	request.CollectorID = &collector.ID

	err = s.repoReq.UpdateRequestPickup(requestId, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update request pickup: %v", err)
	}

	confirmedAt, _ := utils.FormatDateToIndonesianFormat(request.ConfirmedByCollectorAt)

	response := dto.ResponseRequestPickup{
		StatusPickup:           request.StatusPickup,
		CollectorID:            *request.CollectorID,
		ConfirmedByCollectorAt: confirmedAt,
	}

	return &response, nil
}
