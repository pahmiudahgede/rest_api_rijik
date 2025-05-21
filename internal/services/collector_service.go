package services

import (
	"context"
	"errors"

	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
)

type CollectorService interface {
	CreateCollector(ctx context.Context, userID string, req dto.RequestCollectorDTO) error
	AddTrashToCollector(ctx context.Context, collectorID string, req dto.RequestAddAvaibleTrash) error
	GetCollectorByID(ctx context.Context, collectorID string) (*dto.ResponseCollectorDTO, error)
	GetCollectorByUserID(ctx context.Context, userID string) (*dto.ResponseCollectorDTO, error)
	UpdateCollector(ctx context.Context, collectorID string, jobStatus *string, rating float32, addressID string) error
	UpdateAvaibleTrashByCollector(ctx context.Context, collectorID string, updatedTrash []dto.RequestAvaibleTrashbyCollector) error
	DeleteAvaibleTrash(ctx context.Context, trashID string) error
}

type collectorService struct {
	repo      repositories.CollectorRepository
	trashRepo repositories.TrashRepository
}

func NewCollectorService(repo repositories.CollectorRepository, trashRepo repositories.TrashRepository,

) CollectorService {

	return &collectorService{repo: repo, trashRepo: trashRepo}
}

func (s *collectorService) CreateCollector(ctx context.Context, userID string, req dto.RequestCollectorDTO) error {
	collector := &model.Collector{
		UserID:    userID,
		AddressID: req.AddressId,
		JobStatus: "inactive",
		Rating:    5,
	}

	if err := s.repo.CreateCollector(ctx, collector); err != nil {
		return err
	}

	var trashItems []model.AvaibleTrashByCollector
	for _, item := range req.AvaibleTrashbyCollector {
		trashItems = append(trashItems, model.AvaibleTrashByCollector{
			CollectorID:     collector.ID,
			TrashCategoryID: item.TrashId,
			Price:           item.TrashPrice,
		})
	}

	if err := s.repo.AddAvaibleTrash(ctx, trashItems); err != nil {
		return err
	}

	for _, t := range trashItems {
		_ = s.trashRepo.UpdateEstimatedPrice(ctx, t.TrashCategoryID)
	}

	return nil
}

func (s *collectorService) AddTrashToCollector(ctx context.Context, collectorID string, req dto.RequestAddAvaibleTrash) error {
	var trashItems []model.AvaibleTrashByCollector
	for _, item := range req.AvaibleTrash {
		trashItems = append(trashItems, model.AvaibleTrashByCollector{
			CollectorID:     collectorID,
			TrashCategoryID: item.TrashId,
			Price:           item.TrashPrice,
		})
	}
	if err := s.repo.AddAvaibleTrash(ctx, trashItems); err != nil {
		return err
	}

	for _, t := range trashItems {
		_ = s.trashRepo.UpdateEstimatedPrice(ctx, t.TrashCategoryID)
	}

	return nil
}

func (s *collectorService) GetCollectorByID(ctx context.Context, collectorID string) (*dto.ResponseCollectorDTO, error) {
	collector, err := s.repo.GetCollectorByID(ctx, collectorID)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseCollectorDTO{
		ID:        collector.ID,
		UserId:    collector.UserID,
		AddressId: collector.AddressID,
		JobStatus: &collector.JobStatus,
		Rating:    collector.Rating,
		User: &dto.UserResponseDTO{
			ID:    collector.User.ID,
			Name:  collector.User.Name,
			Phone: collector.User.Phone,
		},
		Address: &dto.AddressResponseDTO{
			Province:   collector.Address.Province,
			District:   collector.Address.District,
			Regency:    collector.Address.Regency,
			Village:    collector.Address.Village,
			PostalCode: collector.Address.PostalCode,
			Latitude:   collector.Address.Latitude,
			Longitude:  collector.Address.Longitude,
		},
	}

	for _, item := range collector.AvaibleTrashByCollector {
		response.AvaibleTrashbyCollector = append(response.AvaibleTrashbyCollector, dto.ResponseAvaibleTrashByCollector{
			ID:         item.ID,
			TrashId:    item.TrashCategory.ID,
			TrashName:  item.TrashCategory.Name,
			TrashIcon:  item.TrashCategory.Icon,
			TrashPrice: item.Price,
		})
	}

	return response, nil
}

func (s *collectorService) GetCollectorByUserID(ctx context.Context, userID string) (*dto.ResponseCollectorDTO, error) {
	collector, err := s.repo.GetCollectorByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseCollectorDTO{
		ID:        collector.ID,
		UserId:    collector.UserID,
		AddressId: collector.AddressID,
		JobStatus: &collector.JobStatus,
		Rating:    collector.Rating,
		User: &dto.UserResponseDTO{
			ID:    collector.User.ID,
			Name:  collector.User.Name,
			Phone: collector.User.Phone,
		},
		Address: &dto.AddressResponseDTO{
			Province:   collector.Address.Province,
			District:   collector.Address.District,
			Regency:    collector.Address.Regency,
			Village:    collector.Address.Village,
			PostalCode: collector.Address.PostalCode,
			Latitude:   collector.Address.Latitude,
			Longitude:  collector.Address.Longitude,
		},
	}

	for _, item := range collector.AvaibleTrashByCollector {
		response.AvaibleTrashbyCollector = append(response.AvaibleTrashbyCollector, dto.ResponseAvaibleTrashByCollector{
			ID:         item.ID,
			TrashId:    item.TrashCategory.ID,
			TrashName:  item.TrashCategory.Name,
			TrashIcon:  item.TrashCategory.Icon,
			TrashPrice: item.Price,
		})
	}

	return response, nil
}

func (s *collectorService) UpdateCollector(ctx context.Context, collectorID string, jobStatus *string, rating float32, addressID string) error {
	updates := make(map[string]interface{})

	if jobStatus != nil {
		updates["job_status"] = *jobStatus
	}
	if rating > 0 {
		updates["rating"] = rating
	}
	if addressID != "" {
		updates["address_id"] = addressID
	}

	if len(updates) == 0 {
		return errors.New("tidak ada data yang diubah")
	}

	return s.repo.UpdateCollector(ctx, &model.Collector{ID: collectorID}, updates)
}

func (s *collectorService) UpdateAvaibleTrashByCollector(ctx context.Context, collectorID string, updatedTrash []dto.RequestAvaibleTrashbyCollector) error {
	var updated []model.AvaibleTrashByCollector
	for _, item := range updatedTrash {
		updated = append(updated, model.AvaibleTrashByCollector{
			CollectorID:     collectorID,
			TrashCategoryID: item.TrashId,
			Price:           item.TrashPrice,
		})
	}

	if err := s.repo.UpdateAvaibleTrashByCollector(ctx, collectorID, updated); err != nil {
		return err
	}

	for _, item := range updated {
		_ = s.trashRepo.UpdateEstimatedPrice(ctx, item.TrashCategoryID)
	}

	return nil
}

func (s *collectorService) DeleteAvaibleTrash(ctx context.Context, trashID string) error {
	if trashID == "" {
		return errors.New("trash_id tidak boleh kosong")
	}

	item, err := s.repo.GetTrashItemByID(ctx, trashID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteAvaibleTrash(ctx, trashID); err != nil {
		return err
	}

	return s.trashRepo.UpdateEstimatedPrice(ctx, item.TrashCategoryID)
}
