package repositories

import (
	"context"
	"errors"

	"rijig/config"
	"rijig/model"
)

type CollectorRepository interface {
	CreateCollector(ctx context.Context, collector *model.Collector) error
	AddAvaibleTrash(ctx context.Context, trashItems []model.AvaibleTrashByCollector) error
	GetCollectorByID(ctx context.Context, collectorID string) (*model.Collector, error)
	GetCollectorByUserID(ctx context.Context, userID string) (*model.Collector, error)
	GetTrashItemByID(ctx context.Context, id string) (*model.AvaibleTrashByCollector, error)
	UpdateCollector(ctx context.Context, collector *model.Collector, updates map[string]interface{}) error
	UpdateAvaibleTrashByCollector(ctx context.Context, collectorID string, updatedTrash []model.AvaibleTrashByCollector) error
	DeleteAvaibleTrash(ctx context.Context, trashID string) error

	GetActiveCollectorsWithTrashAndAddress(ctx context.Context) ([]model.Collector, error)
	GetCollectorWithAddressAndTrash(ctx context.Context, collectorID string) (*model.Collector, error)
}

type collectorRepository struct {
}

func NewCollectorRepository() CollectorRepository {
	return &collectorRepository{}
}

func (r *collectorRepository) CreateCollector(ctx context.Context, collector *model.Collector) error {
	return config.DB.WithContext(ctx).Create(collector).Error
}

func (r *collectorRepository) AddAvaibleTrash(ctx context.Context, trashItems []model.AvaibleTrashByCollector) error {
	if len(trashItems) == 0 {
		return nil
	}
	return config.DB.WithContext(ctx).Create(&trashItems).Error
}

func (r *collectorRepository) GetCollectorByID(ctx context.Context, collectorID string) (*model.Collector, error) {
	var collector model.Collector
	err := config.DB.WithContext(ctx).
		Preload("User").
		Preload("Address").
		Preload("AvaibleTrashByCollector.TrashCategory").
		First(&collector, "id = ?", collectorID).Error

	if err != nil {
		return nil, err
	}
	return &collector, nil
}

func (r *collectorRepository) GetCollectorByUserID(ctx context.Context, userID string) (*model.Collector, error) {
	var collector model.Collector
	err := config.DB.WithContext(ctx).
		Preload("User").
		Preload("Address").
		Preload("AvaibleTrashByCollector.TrashCategory").
		First(&collector, "user_id = ?", userID).Error

	if err != nil {
		return nil, err
	}
	return &collector, nil
}

func (r *collectorRepository) GetTrashItemByID(ctx context.Context, id string) (*model.AvaibleTrashByCollector, error) {
	var item model.AvaibleTrashByCollector
	if err := config.DB.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *collectorRepository) UpdateCollector(ctx context.Context, collector *model.Collector, updates map[string]interface{}) error {
	return config.DB.WithContext(ctx).
		Model(&model.Collector{}).
		Where("id = ?", collector.ID).
		Updates(updates).Error
}

func (r *collectorRepository) UpdateAvaibleTrashByCollector(ctx context.Context, collectorID string, updatedTrash []model.AvaibleTrashByCollector) error {
	for _, trash := range updatedTrash {
		err := config.DB.WithContext(ctx).
			Model(&model.AvaibleTrashByCollector{}).
			Where("collector_id = ? AND trash_category_id = ?", collectorID, trash.TrashCategoryID).
			Update("price", trash.Price).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *collectorRepository) DeleteAvaibleTrash(ctx context.Context, trashID string) error {
	if trashID == "" {
		return errors.New("trash_id cannot be empty")
	}
	return config.DB.WithContext(ctx).
		Delete(&model.AvaibleTrashByCollector{}, "id = ?", trashID).Error
}


// 
func (r *collectorRepository) GetActiveCollectorsWithTrashAndAddress(ctx context.Context) ([]model.Collector, error) {
	var collectors []model.Collector
	err := config.DB.WithContext(ctx).
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
	err := config.DB.WithContext(ctx).
		Preload("Address").
		Preload("AvaibleTrashbyCollector").
		Where("id = ?", collectorID).
		First(&collector).Error

	if err != nil {
		return nil, err
	}
	return &collector, nil
}
