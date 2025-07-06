package requestpickup

import (
	"context"
	"fmt"
	"rijig/internal/collector"
	"rijig/utils"
)

type PickupMatchingService interface {
	FindNearbyCollectorsForPickup(ctx context.Context, pickupID string) ([]collector.NearbyCollectorDTO, error)
	FindAvailableRequestsForCollector(ctx context.Context, collectorID string) ([]PickupRequestForCollectorDTO, error)
}

type pickupMatchingService struct {
	pickupRepo    RequestPickupRepository
	collectorRepo collector.CollectorRepository
}

func NewPickupMatchingService(pickupRepo RequestPickupRepository,
	collectorRepo collector.CollectorRepository) PickupMatchingService {
	return &pickupMatchingService{
		pickupRepo:    pickupRepo,
		collectorRepo: collectorRepo,
	}
}

func (s *pickupMatchingService) FindNearbyCollectorsForPickup(ctx context.Context, pickupID string) ([]collector.NearbyCollectorDTO, error) {
	pickup, err := s.pickupRepo.GetPickupWithItemsAndAddress(ctx, pickupID)
	if err != nil {
		return nil, fmt.Errorf("pickup tidak ditemukan: %w", err)
	}

	userCoord := utils.Coord{
		Lat: pickup.Address.Latitude,
		Lon: pickup.Address.Longitude,
	}

	requestedTrash := make(map[string]bool)
	for _, item := range pickup.RequestItems {
		requestedTrash[item.TrashCategoryId] = true
	}

	collectors, err := s.collectorRepo.GetActiveCollectorsWithTrashAndAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data collector: %w", err)
	}

	var result []collector.NearbyCollectorDTO
	for _, col := range collectors {
		coord := utils.Coord{
			Lat: col.Address.Latitude,
			Lon: col.Address.Longitude,
		}

		_, km := utils.Distance(userCoord, coord)
		if km > 10 {
			continue
		}

		var matchedTrash []string
		for _, item := range col.AvaibleTrashByCollector {
			if requestedTrash[item.TrashCategoryID] {
				matchedTrash = append(matchedTrash, item.TrashCategoryID)
			}
		}

		if len(matchedTrash) == 0 {
			continue
		}

		result = append(result, collector.NearbyCollectorDTO{
			CollectorID:  col.ID,
			Name:         col.User.Name,
			Phone:        col.User.Phone,
			Rating:       col.Rating,
			Latitude:     col.Address.Latitude,
			Longitude:    col.Address.Longitude,
			DistanceKm:   km,
			MatchedTrash: matchedTrash,
		})
	}

	return result, nil
}

// terdpaat error seperti ini: "undefined: dto.PickupRequestForCollectorDTO" dan seprti ini: s.collectorRepo.GetCollectorWithAddressAndTrash undefined (type repositories.CollectorRepository has no field or method GetCollectorWithAddressAndTrash) pada kode berikut:

func (s *pickupMatchingService) FindAvailableRequestsForCollector(ctx context.Context, collectorID string) ([]PickupRequestForCollectorDTO, error) {
	collector, err := s.collectorRepo.GetCollectorWithAddressAndTrash(ctx, collectorID)
	if err != nil {
		return nil, fmt.Errorf("collector tidak ditemukan: %w", err)
	}

	pickupList, err := s.pickupRepo.GetAllAutomaticRequestsWithAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pickup otomatis: %w", err)
	}

	collectorCoord := utils.Coord{
		Lat: collector.Address.Latitude,
		Lon: collector.Address.Longitude,
	}

	// map trash collector
	collectorTrash := make(map[string]bool)
	for _, t := range collector.AvaibleTrashByCollector {
		collectorTrash[t.TrashCategoryID] = true
	}

	var results []PickupRequestForCollectorDTO
	for _, p := range pickupList {
		if p.StatusPickup != "waiting_collector" {
			continue
		}
		coord := utils.Coord{
			Lat: p.Address.Latitude,
			Lon: p.Address.Longitude,
		}
		_, km := utils.Distance(collectorCoord, coord)
		if km > 10 {
			continue
		}

		match := false
		var matchedTrash []string
		for _, item := range p.RequestItems {
			if collectorTrash[item.TrashCategoryId] {
				match = true
				matchedTrash = append(matchedTrash, item.TrashCategoryId)
			}
		}
		if match {
			results = append(results, PickupRequestForCollectorDTO{
				PickupID:     p.ID,
				UserID:       p.UserId,
				Latitude:     p.Address.Latitude,
				Longitude:    p.Address.Longitude,
				DistanceKm:   km,
				MatchedTrash: matchedTrash,
			})
		}
	}

	return results, nil
}
