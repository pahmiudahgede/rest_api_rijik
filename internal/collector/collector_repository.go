package collector

import (
	"context"
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type CollectorRepository interface {
	Create(ctx context.Context, collector *model.Collector) error
	GetByID(ctx context.Context, id string) (*model.Collector, error)
	GetByUserID(ctx context.Context, userID string) (*model.Collector, error)
	Update(ctx context.Context, collector *model.Collector) error
	Delete(ctx context.Context, UserID string) error
	List(ctx context.Context, limit, offset int) ([]*model.Collector, int64, error)

	GetActiveCollectors(ctx context.Context, limit, offset int) ([]*model.Collector, int64, error)
	GetCollectorsByAddress(ctx context.Context, addressID string, limit, offset int) ([]*model.Collector, int64, error)
	GetCollectorsByTrashCategory(ctx context.Context, trashCategoryID string, limit, offset int) ([]*model.Collector, int64, error)
	UpdateJobStatus(ctx context.Context, id string, jobStatus string) error
	UpdateRating(ctx context.Context, id string, rating float32) error

	CreateAvailableTrash(ctx context.Context, availableTrash *model.AvaibleTrashByCollector) error
	GetAvailableTrashByCollectorID(ctx context.Context, collectorID string) ([]*model.AvaibleTrashByCollector, error)
	UpdateAvailableTrash(ctx context.Context, availableTrash *model.AvaibleTrashByCollector) error
	DeleteAvailableTrash(ctx context.Context, id string) error
	BulkCreateAvailableTrash(ctx context.Context, availableTrashList []*model.AvaibleTrashByCollector) error
	BulkUpdateAvailableTrash(ctx context.Context, collectorID string, availableTrashList []*model.AvaibleTrashByCollector) error
	DeleteAvailableTrashByCollectorID(ctx context.Context, collectorID string) error

	GetActiveCollectorsWithTrashAndAddress(ctx context.Context) ([]model.Collector, error)
	GetCollectorWithAddressAndTrash(ctx context.Context, collectorID string) (*model.Collector, error)

	WithTx(tx *gorm.DB) CollectorRepository
}

type collectorRepository struct {
	db *gorm.DB
}

func NewCollectorRepository(db *gorm.DB) CollectorRepository {
	return &collectorRepository{
		db: db,
	}
}

func (r *collectorRepository) WithTx(tx *gorm.DB) CollectorRepository {
	return &collectorRepository{
		db: tx,
	}
}

func (r *collectorRepository) Create(ctx context.Context, collector *model.Collector) error {
	if err := r.db.WithContext(ctx).Create(collector).Error; err != nil {
		return fmt.Errorf("failed to create collector: %w", err)
	}
	return nil
}

func (r *collectorRepository) GetActiveCollectorsWithTrashAndAddress(ctx context.Context) ([]model.Collector, error) {
	var collectors []model.Collector
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Address").
		Preload("AvaibleTrashbyCollector.TrashCategory").
		Where("job_status = ?", "active").
		Find(&collectors).Error

	if err != nil {
		return nil, err
	}

	return collectors, nil
}

func (r *collectorRepository) GetCollectorWithAddressAndTrash(ctx context.Context, collectorID string) (*model.Collector, error) {
	var collector model.Collector
	err := r.db.WithContext(ctx).
		Preload("Address").
		Preload("AvaibleTrashbyCollector").
		Where("id = ?", collectorID).
		First(&collector).Error

	if err != nil {
		return nil, err
	}
	return &collector, nil
}

func (r *collectorRepository) GetByID(ctx context.Context, id string) (*model.Collector, error) {
	var collector model.Collector

	err := r.db.WithContext(ctx).
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Where("id = ?", id).
		First(&collector).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("collector with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get collector by id: %w", err)
	}

	return &collector, nil
}

func (r *collectorRepository) GetByUserID(ctx context.Context, userID string) (*model.Collector, error) {
	var collector model.Collector

	err := r.db.WithContext(ctx).
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Where("user_id = ?", userID).
		First(&collector).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("collector with user_id %s not found", userID)
		}
		return nil, fmt.Errorf("failed to get collector by user_id: %w", err)
	}

	return &collector, nil
}

func (r *collectorRepository) Update(ctx context.Context, collector *model.Collector) error {
	if err := r.db.WithContext(ctx).Save(collector).Error; err != nil {
		return fmt.Errorf("failed to update collector: %w", err)
	}
	return nil
}

func (r *collectorRepository) Delete(ctx context.Context, UserID string) error {
	result := r.db.WithContext(ctx).Delete(&model.Collector{}, "user_id = ?", UserID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete collector: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("collector with user_id %s not found", UserID)
	}

	return nil
}

func (r *collectorRepository) List(ctx context.Context, limit, offset int) ([]*model.Collector, int64, error) {
	var collectors []*model.Collector
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Collector{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count collectors: %w", err)
	}

	err := r.db.WithContext(ctx).
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Limit(limit).
		Offset(offset).
		Find(&collectors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list collectors: %w", err)
	}

	return collectors, total, nil
}

func (r *collectorRepository) GetActiveCollectors(ctx context.Context, limit, offset int) ([]*model.Collector, int64, error) {
	var collectors []*model.Collector
	var total int64

	query := r.db.WithContext(ctx).Where("job_status = ?", "active")

	if err := query.Model(&model.Collector{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count active collectors: %w", err)
	}

	err := query.
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Limit(limit).
		Offset(offset).
		Find(&collectors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active collectors: %w", err)
	}

	return collectors, total, nil
}

func (r *collectorRepository) GetCollectorsByAddress(ctx context.Context, addressID string, limit, offset int) ([]*model.Collector, int64, error) {
	var collectors []*model.Collector
	var total int64

	query := r.db.WithContext(ctx).Where("address_id = ?", addressID)

	if err := query.Model(&model.Collector{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count collectors by address: %w", err)
	}

	err := query.
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Limit(limit).
		Offset(offset).
		Find(&collectors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collectors by address: %w", err)
	}

	return collectors, total, nil
}

func (r *collectorRepository) GetCollectorsByTrashCategory(ctx context.Context, trashCategoryID string, limit, offset int) ([]*model.Collector, int64, error) {
	var collectors []*model.Collector
	var total int64

	subQuery := r.db.WithContext(ctx).
		Table("avaible_trash_by_collectors").
		Select("collector_id").
		Where("trash_category_id = ?", trashCategoryID)

	query := r.db.WithContext(ctx).
		Where("id IN (?)", subQuery)

	if err := query.Model(&model.Collector{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count collectors by trash category: %w", err)
	}

	err := query.
		Preload("Address").
		Preload("AvaibleTrashByCollector").
		Preload("AvaibleTrashByCollector.TrashCategory").
		Limit(limit).
		Offset(offset).
		Find(&collectors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collectors by trash category: %w", err)
	}

	return collectors, total, nil
}

func (r *collectorRepository) UpdateJobStatus(ctx context.Context, id string, jobStatus string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Collector{}).
		Where("id = ?", id).
		Update("job_status", jobStatus)

	if result.Error != nil {
		return fmt.Errorf("failed to update job status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("collector with id %s not found", id)
	}

	return nil
}

func (r *collectorRepository) UpdateRating(ctx context.Context, id string, rating float32) error {
	result := r.db.WithContext(ctx).
		Model(&model.Collector{}).
		Where("id = ?", id).
		Update("rating", rating)

	if result.Error != nil {
		return fmt.Errorf("failed to update rating: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("collector with id %s not found", id)
	}

	return nil
}

func (r *collectorRepository) CreateAvailableTrash(ctx context.Context, availableTrash *model.AvaibleTrashByCollector) error {
	if err := r.db.WithContext(ctx).Create(availableTrash).Error; err != nil {
		return fmt.Errorf("failed to create available trash: %w", err)
	}
	return nil
}

func (r *collectorRepository) GetAvailableTrashByCollectorID(ctx context.Context, collectorID string) ([]*model.AvaibleTrashByCollector, error) {
	var availableTrash []*model.AvaibleTrashByCollector

	err := r.db.WithContext(ctx).
		Preload("TrashCategory").
		Where("collector_id = ?", collectorID).
		Find(&availableTrash).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available trash by collector id: %w", err)
	}

	return availableTrash, nil
}

func (r *collectorRepository) UpdateAvailableTrash(ctx context.Context, availableTrash *model.AvaibleTrashByCollector) error {
	if err := r.db.WithContext(ctx).Save(availableTrash).Error; err != nil {
		return fmt.Errorf("failed to update available trash: %w", err)
	}
	return nil
}

func (r *collectorRepository) DeleteAvailableTrash(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.AvaibleTrashByCollector{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete available trash: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("available trash with id %s not found", id)
	}

	return nil
}

func (r *collectorRepository) BulkCreateAvailableTrash(ctx context.Context, availableTrashList []*model.AvaibleTrashByCollector) error {
	if len(availableTrashList) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(availableTrashList, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk create available trash: %w", err)
	}

	return nil
}

func (r *collectorRepository) BulkUpdateAvailableTrash(ctx context.Context, collectorID string, availableTrashList []*model.AvaibleTrashByCollector) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("collector_id = ?", collectorID).Delete(&model.AvaibleTrashByCollector{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing available trash: %w", err)
		}

		if len(availableTrashList) > 0 {
			for _, item := range availableTrashList {
				item.CollectorID = collectorID
			}

			if err := tx.CreateInBatches(availableTrashList, 100).Error; err != nil {
				return fmt.Errorf("failed to create new available trash: %w", err)
			}
		}

		return nil
	})
}

func (r *collectorRepository) DeleteAvailableTrashByCollectorID(ctx context.Context, collectorID string) error {
	if err := r.db.WithContext(ctx).Where("collector_id = ?", collectorID).Delete(&model.AvaibleTrashByCollector{}).Error; err != nil {
		return fmt.Errorf("failed to delete available trash by collector id: %w", err)
	}
	return nil
}
