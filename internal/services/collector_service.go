package services

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/utils"
	"time"
)

type CollectorService interface {
	FindCollectorsNearby(userId string) ([]dto.ResponseCollectorDTO, error)
	ConfirmRequestPickup(requestId, collectorId string) (*dto.ResponseRequestPickup, error)
	ConfirmRequestManualPickup(requestId, collectorId string) (any, error)
}

type collectorService struct {
	repo        repositories.CollectorRepository
	repoColl    repositories.RequestPickupRepository
	repoAddress repositories.AddressRepository
	repoUser    repositories.UserProfilRepository
}

func NewCollectorService(repo repositories.CollectorRepository,
	repoColl repositories.RequestPickupRepository,
	repoAddress repositories.AddressRepository,
	repoUser repositories.UserProfilRepository) CollectorService {
	return &collectorService{repo: repo, repoColl: repoColl, repoAddress: repoAddress, repoUser: repoUser}
}

func (s *collectorService) FindCollectorsNearby(userId string) ([]dto.ResponseCollectorDTO, error) {
	collectors, err := s.repo.FindActiveCollectors()
	if err != nil {
		return nil, fmt.Errorf("error fetching active collectors: %v", err)
	}

	var avaibleCollectResp []dto.ResponseCollectorDTO

	for _, collector := range collectors {

		request, err := s.repoColl.FindRequestPickupByAddressAndStatus(userId, "waiting_collector", "otomatis")
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan data request pickup dengan userid: %v", err)
		}

		_, distance := utils.Distance(
			utils.Coord{Lat: request.Address.Latitude, Lon: request.Address.Longitude},
			utils.Coord{Lat: collector.Address.Latitude, Lon: collector.Address.Longitude},
		)

		if distance <= 20 {

			mappedRequest := dto.ResponseCollectorDTO{
				ID:        collector.ID,
				UserId:    collector.UserID,
				AddressId: collector.AddressId,

				Rating: collector.Rating,
			}

			user, err := s.repoUser.FindByID(collector.UserID)
			if err != nil {
				return nil, fmt.Errorf("error fetching user data: %v", err)
			}
			mappedRequest.User = []dto.UserResponseDTO{
				{
					Name:  user.Name,
					Phone: user.Phone,
				},
			}

			address, err := s.repoAddress.FindAddressByID(collector.AddressId)
			if err != nil {
				return nil, fmt.Errorf("error fetching address data: %v", err)
			}
			mappedRequest.Address = []dto.AddressResponseDTO{
				{
					District: address.District,
					Village:  address.Village,
					Detail:   address.Detail,
				},
			}

			avaibleCollectResp = append(avaibleCollectResp, mappedRequest)
		}
	}

	return avaibleCollectResp, nil
}

func (s *collectorService) ConfirmRequestPickup(requestId, collectorId string) (*dto.ResponseRequestPickup, error) {
	request, err := s.repoColl.FindRequestPickupByID(requestId)
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
	*request.ConfirmedByCollectorAt = time.Now()

	err = s.repoColl.UpdateRequestPickup(requestId, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update request pickup: %v", err)
	}

	confirmedAt, _ := utils.FormatDateToIndonesianFormat(*request.ConfirmedByCollectorAt)

	response := dto.ResponseRequestPickup{
		StatusPickup:           request.StatusPickup,
		CollectorID:            *request.CollectorID,
		ConfirmedByCollectorAt: confirmedAt,
	}

	return &response, nil
}

func (s *collectorService) ConfirmRequestManualPickup(requestId, collectorId string) (any, error) {

	request, err := s.repoColl.FindRequestPickupByID(requestId)
	if err != nil {
		return nil, fmt.Errorf("collector tidak ditemukan: %v", err)
	}

	coll, err := s.repo.FindCollectorByIdWithoutAddr(collectorId)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if coll.ID != *request.CollectorID {
		return nil, fmt.Errorf("collectorid tidak sesuai dengan request")
	}

	request.StatusPickup = "confirmed"
	*request.ConfirmedByCollectorAt = time.Now()

	err = s.repoColl.UpdateRequestPickup(requestId, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update request pickup: %v", err)
	}

	return "berhasil konfirmasi request pickup", nil
}