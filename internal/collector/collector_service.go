package collector

import (
	"context"
	"fmt"
	"rijig/internal/address"
	"rijig/internal/trash"
	"rijig/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CollectorService interface {
	CreateCollector(ctx context.Context, req *CreateCollectorRequest, UserID string) (*CollectorResponse, error)
	GetCollectorByID(ctx context.Context, id string) (*CollectorResponse, error)
	GetCollectorByUserID(ctx context.Context, userID string) (*CollectorResponse, error)
	UpdateCollector(ctx context.Context, UserID string, req *UpdateCollectorRequest) (*CollectorResponse, error)
	DeleteCollector(ctx context.Context, UserID string) error
	ListCollectors(ctx context.Context, limit, offset int) ([]*CollectorResponse, int64, error)

	GetActiveCollectors(ctx context.Context, limit, offset int) ([]*CollectorResponse, int64, error)
	GetCollectorsByAddress(ctx context.Context, addressID string, limit, offset int) ([]*CollectorResponse, int64, error)
	GetCollectorsByTrashCategory(ctx context.Context, trashCategoryID string, limit, offset int) ([]*CollectorResponse, int64, error)

	UpdateJobStatus(ctx context.Context, id string, jobStatus string) error
	UpdateRating(ctx context.Context, id string, rating float32) error
	UpdateAvailableTrash(ctx context.Context, collectorID string, availableTrashItems []CreateAvailableTrashRequest) error
}

type collectorService struct {
	collectorRepo CollectorRepository
	db            *gorm.DB
}

func NewCollectorService(collectorRepo CollectorRepository, db *gorm.DB) CollectorService {
	return &collectorService{
		collectorRepo: collectorRepo,
		db:            db,
	}
}

func (s *collectorService) CreateCollector(ctx context.Context, req *CreateCollectorRequest, UserID string) (*CollectorResponse, error) {

	existingCollector, err := s.collectorRepo.GetByUserID(ctx, UserID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, fmt.Errorf("failed to check existing collector: %w", err)
	}
	if existingCollector != nil {
		return nil, fmt.Errorf("collector already exists for user_id: %s", req.UserID)
	}

	collector := &model.Collector{
		UserID:    UserID,
		JobStatus: "inactive",
		AddressID: req.AddressID,
		Rating:    5.0,
	}

	if req.JobStatus != "" {
		collector.JobStatus = req.JobStatus
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		collectorRepoTx := s.collectorRepo.WithTx(tx)

		if err := collectorRepoTx.Create(ctx, collector); err != nil {
			return fmt.Errorf("failed to create collector: %w", err)
		}

		if len(req.AvailableTrashItems) > 0 {
			availableTrashList := s.buildAvailableTrashList(collector.ID, req.AvailableTrashItems)
			if err := collectorRepoTx.BulkCreateAvailableTrash(ctx, availableTrashList); err != nil {
				return fmt.Errorf("failed to create available trash items: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	createdCollector, err := s.collectorRepo.GetByID(ctx, collector.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created collector: %w", err)
	}

	return s.toCollectorResponse(createdCollector), nil
}

func (s *collectorService) GetCollectorByID(ctx context.Context, id string) (*CollectorResponse, error) {
	collector, err := s.collectorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCollectorResponse(collector), nil
}

func (s *collectorService) GetCollectorByUserID(ctx context.Context, userID string) (*CollectorResponse, error) {
	collector, err := s.collectorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toCollectorResponse(collector), nil
}

func (s *collectorService) UpdateCollector(ctx context.Context, UserID string, req *UpdateCollectorRequest) (*CollectorResponse, error) {

	collector, err := s.collectorRepo.GetByUserID(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collector: %w", err)
	}

	needsUpdate := s.checkCollectorNeedsUpdate(collector, req)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		collectorRepoTx := s.collectorRepo.WithTx(tx)

		if needsUpdate {
			s.applyCollectorUpdates(collector, req)
			collector.UpdatedAt = time.Now()

			if err := collectorRepoTx.Update(ctx, collector); err != nil {
				return fmt.Errorf("failed to update collector: %w", err)
			}
		}

		if len(req.AvailableTrashItems) > 0 {
			availableTrashList := s.buildAvailableTrashList(collector.ID, req.AvailableTrashItems)
			if err := collectorRepoTx.BulkUpdateAvailableTrash(ctx, collector.ID, availableTrashList); err != nil {
				return fmt.Errorf("failed to update available trash items: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedCollector, err := s.collectorRepo.GetByUserID(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated collector: %w", err)
	}

	return s.toCollectorResponse(updatedCollector), nil
}

func (s *collectorService) DeleteCollector(ctx context.Context, UserID string) error {

	_, err := s.collectorRepo.GetByUserID(ctx, UserID)
	if err != nil {
		return fmt.Errorf("collector not found: %w", err)
	}

	if err := s.collectorRepo.Delete(ctx, UserID); err != nil {
		return fmt.Errorf("failed to delete collector: %w", err)
	}

	return nil
}

func (s *collectorService) ListCollectors(ctx context.Context, limit, offset int) ([]*CollectorResponse, int64, error) {

	limit, offset = s.normalizePagination(limit, offset)

	collectors, total, err := s.collectorRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list collectors: %w", err)
	}

	return s.buildCollectorResponseList(collectors), total, nil
}

func (s *collectorService) GetActiveCollectors(ctx context.Context, limit, offset int) ([]*CollectorResponse, int64, error) {

	limit, offset = s.normalizePagination(limit, offset)

	collectors, total, err := s.collectorRepo.GetActiveCollectors(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active collectors: %w", err)
	}

	return s.buildCollectorResponseList(collectors), total, nil
}

func (s *collectorService) GetCollectorsByAddress(ctx context.Context, addressID string, limit, offset int) ([]*CollectorResponse, int64, error) {

	limit, offset = s.normalizePagination(limit, offset)

	collectors, total, err := s.collectorRepo.GetCollectorsByAddress(ctx, addressID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collectors by address: %w", err)
	}

	return s.buildCollectorResponseList(collectors), total, nil
}

func (s *collectorService) GetCollectorsByTrashCategory(ctx context.Context, trashCategoryID string, limit, offset int) ([]*CollectorResponse, int64, error) {

	limit, offset = s.normalizePagination(limit, offset)

	collectors, total, err := s.collectorRepo.GetCollectorsByTrashCategory(ctx, trashCategoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collectors by trash category: %w", err)
	}

	return s.buildCollectorResponseList(collectors), total, nil
}

func (s *collectorService) UpdateJobStatus(ctx context.Context, id string, jobStatus string) error {
	if err := s.collectorRepo.UpdateJobStatus(ctx, id, jobStatus); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

func (s *collectorService) UpdateRating(ctx context.Context, id string, rating float32) error {
	if err := s.collectorRepo.UpdateRating(ctx, id, rating); err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	return nil
}

func (s *collectorService) UpdateAvailableTrash(ctx context.Context, collectorID string, availableTrashItems []CreateAvailableTrashRequest) error {
	availableTrashList := s.buildAvailableTrashList(collectorID, availableTrashItems)

	if err := s.collectorRepo.BulkUpdateAvailableTrash(ctx, collectorID, availableTrashList); err != nil {
		return fmt.Errorf("failed to update available trash: %w", err)
	}

	return nil
}

func (s *collectorService) buildAvailableTrashList(collectorID string, items []CreateAvailableTrashRequest) []*model.AvaibleTrashByCollector {
	availableTrashList := make([]*model.AvaibleTrashByCollector, 0, len(items))
	for _, item := range items {
		availableTrash := &model.AvaibleTrashByCollector{
			CollectorID:     collectorID,
			TrashCategoryID: item.TrashCategoryID,
			Price:           item.Price,
		}
		availableTrashList = append(availableTrashList, availableTrash)
	}
	return availableTrashList
}

func (s *collectorService) checkCollectorNeedsUpdate(collector *model.Collector, req *UpdateCollectorRequest) bool {
	if req.JobStatus != "" && req.JobStatus != collector.JobStatus {
		return true
	}
	if req.AddressID != "" && req.AddressID != collector.AddressID {
		return true
	}
	return false
}

func (s *collectorService) applyCollectorUpdates(collector *model.Collector, req *UpdateCollectorRequest) {
	if req.JobStatus != "" {
		collector.JobStatus = req.JobStatus
	}
	if req.AddressID != "" {
		collector.AddressID = req.AddressID
	}
}

func (s *collectorService) normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if limit > 100 {
		limit = 100
	}
	return limit, offset
}

func (s *collectorService) buildCollectorResponseList(collectors []*model.Collector) []*CollectorResponse {
	responses := make([]*CollectorResponse, 0, len(collectors))
	for _, collector := range collectors {
		responses = append(responses, s.toCollectorResponse(collector))
	}
	return responses
}

func (s *collectorService) toCollectorResponse(collector *model.Collector) *CollectorResponse {
	response := &CollectorResponse{
		ID:             collector.ID,
		UserID:         collector.UserID,
		JobStatus:      collector.JobStatus,
		Rating:         collector.Rating,
		AddressID:      collector.AddressID,
		AvailableTrash: make([]AvailableTrashResponse, 0),
		CreatedAt:      collector.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      collector.UpdatedAt.Format(time.RFC3339),
	}

	if collector.Address.ID != "" {
		response.Address = &address.AddressResponseDTO{
			ID:         collector.Address.ID,
			UserID:     collector.Address.UserID,
			Province:   collector.Address.Province,
			Regency:    collector.Address.Regency,
			District:   collector.Address.District,
			Village:    collector.Address.Village,
			PostalCode: collector.Address.PostalCode,
			Detail:     collector.Address.Detail,
			Latitude:   collector.Address.Latitude,
			Longitude:  collector.Address.Longitude,
			CreatedAt:  collector.Address.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  collector.Address.UpdatedAt.Format(time.RFC3339),
		}
	}

	for _, availableTrash := range collector.AvaibleTrashByCollector {
		trashResponse := AvailableTrashResponse{
			ID:              availableTrash.ID,
			CollectorID:     availableTrash.CollectorID,
			TrashCategoryID: availableTrash.TrashCategoryID,
			Price:           availableTrash.Price,
		}

		if availableTrash.TrashCategory.ID != "" {
			trashResponse.TrashCategory = &trash.ResponseTrashCategoryDTO{
				ID:             availableTrash.TrashCategory.ID,
				TrashName:      availableTrash.TrashCategory.Name,
				TrashIcon:      availableTrash.TrashCategory.IconTrash,
				EstimatedPrice: availableTrash.TrashCategory.EstimatedPrice,
				Variety:        availableTrash.TrashCategory.Variety,
				CreatedAt:      availableTrash.TrashCategory.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      availableTrash.TrashCategory.UpdatedAt.Format(time.RFC3339),
			}
		}

		response.AvailableTrash = append(response.AvailableTrash, trashResponse)
	}

	return response
}
