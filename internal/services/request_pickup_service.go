package services

import (
	"context"
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"time"
)

type RequestPickupService interface {
	ConvertCartToRequestPickup(ctx context.Context, userID string, req dto.RequestPickupDTO) error
	AssignCollectorToRequest(ctx context.Context, pickupID string, req dto.SelectCollectorDTO) error
	FindRequestsAssignedToCollector(ctx context.Context, collectorID string) ([]dto.AssignedPickupDTO, error)
	ConfirmPickupByCollector(ctx context.Context, pickupID string, confirmedAt time.Time) error
	UpdatePickupStatusToPickingUp(ctx context.Context, pickupID string) error
	UpdateActualPickupItems(ctx context.Context, pickupID string, items []dto.UpdateRequestPickupItemDTO) error
}

type requestPickupService struct {
	trashRepo   repositories.TrashRepository
	pickupRepo  repositories.RequestPickupRepository
	cartService CartService
	historyRepo repositories.PickupStatusHistoryRepository
}

func NewRequestPickupService(trashRepo repositories.TrashRepository, pickupRepo repositories.RequestPickupRepository, cartService CartService, historyRepo repositories.PickupStatusHistoryRepository) RequestPickupService {
	return &requestPickupService{
		trashRepo:   trashRepo,
		pickupRepo:  pickupRepo,
		cartService: cartService,
		historyRepo: historyRepo,
	}
}

func (s *requestPickupService) ConvertCartToRequestPickup(ctx context.Context, userID string, req dto.RequestPickupDTO) error {
	cart, err := s.cartService.GetCart(ctx, userID)
	if err != nil || len(cart.CartItems) == 0 {
		return fmt.Errorf("cart kosong atau tidak ditemukan")
	}

	var requestItems []model.RequestPickupItem
	for _, item := range cart.CartItems {
		trash, err := s.trashRepo.GetTrashCategoryByID(ctx, item.TrashID)
		if err != nil {
			continue
		}
		subtotal := float64(item.Amount) * trash.EstimatedPrice

		requestItems = append(requestItems, model.RequestPickupItem{
			TrashCategoryId:        item.TrashID,
			EstimatedAmount:        float64(item.Amount),
			EstimatedPricePerKg:    trash.EstimatedPrice,
			EstimatedSubtotalPrice: subtotal,
		})
	}

	if len(requestItems) == 0 {
		return fmt.Errorf("tidak ada item valid dalam cart")
	}

	pickup := model.RequestPickup{
		UserId:        userID,
		AddressId:     req.AddressID,
		RequestMethod: req.RequestMethod,
		Notes:         req.Notes,
		StatusPickup:  "waiting_collector",
		RequestItems:  requestItems,
	}

	if err := s.pickupRepo.CreateRequestPickup(ctx, &pickup); err != nil {
		return fmt.Errorf("gagal menyimpan request pickup: %w", err)
	}

	if err := s.cartService.ClearCart(ctx, userID); err != nil {
		return fmt.Errorf("request berhasil, tapi gagal hapus cart: %w", err)
	}

	return nil
}

func (s *requestPickupService) AssignCollectorToRequest(ctx context.Context, pickupID string, req dto.SelectCollectorDTO) error {
	if req.CollectorID == "" {
		return fmt.Errorf("collector_id tidak boleh kosong")
	}
	return s.pickupRepo.UpdateCollectorID(ctx, pickupID, req.CollectorID)
}

func (s *requestPickupService) FindRequestsAssignedToCollector(ctx context.Context, collectorID string) ([]dto.AssignedPickupDTO, error) {
	pickups, err := s.pickupRepo.GetRequestsAssignedToCollector(ctx, collectorID)
	if err != nil {
		return nil, err
	}

	var result []dto.AssignedPickupDTO
	for _, p := range pickups {
		var matchedTrash []string
		for _, item := range p.RequestItems {
			matchedTrash = append(matchedTrash, item.TrashCategoryId)
		}

		result = append(result, dto.AssignedPickupDTO{
			PickupID:     p.ID,
			UserID:       p.UserId,
			UserName:     p.User.Name,
			Latitude:     p.Address.Latitude,
			Longitude:    p.Address.Longitude,
			Notes:        p.Notes,
			MatchedTrash: matchedTrash,
		})
	}

	return result, nil
}

func (s *requestPickupService) ConfirmPickupByCollector(ctx context.Context, pickupID string, confirmedAt time.Time) error {
	return s.pickupRepo.UpdatePickupStatusAndConfirmationTime(ctx, pickupID, "confirmed_by_collector", confirmedAt)
}

func (s *requestPickupService) UpdatePickupStatusToPickingUp(ctx context.Context, pickupID string) error {
	err := s.pickupRepo.UpdatePickupStatus(ctx, pickupID, "collector_are_picking_up")
	if err != nil {
		return err
	}
	return s.historyRepo.CreateStatusHistory(ctx, model.PickupStatusHistory{
		RequestID:     pickupID,
		Status:        "collector_are_picking_up",
		ChangedAt:     time.Now(),
		ChangedByID:   "collector",
		ChangedByRole: "collector",
	})
}

func (s *requestPickupService) UpdateActualPickupItems(ctx context.Context, pickupID string, items []dto.UpdateRequestPickupItemDTO) error {
	return s.pickupRepo.UpdateRequestPickupItemsAmountAndPrice(ctx, pickupID, items)
}
