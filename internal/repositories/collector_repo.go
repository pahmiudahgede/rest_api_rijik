package repositories

import (
	"context"
	"errors"

	// "fmt"

	// "log"
	"rijig/config"
	"rijig/model"
	// "gorm.io/gorm"
)

type CollectorRepository interface {
	// FindActiveCollectors() ([]model.Collector, error)
	// FindCollectorById(collector_id string) (*model.Collector, error)
	// FindCollectorByIdWithoutAddr(collector_id string) (*model.Collector, error)
	// CreateCollector(collector *model.Collector) error
	// UpdateCollector(userId string, jobStatus string) (*model.Collector, error)

	CreateCollector(ctx context.Context, collector *model.Collector) error
	AddAvaibleTrash(ctx context.Context, trashItems []model.AvaibleTrashByCollector) error
	GetCollectorByID(ctx context.Context, collectorID string) (*model.Collector, error)
	GetCollectorByUserID(ctx context.Context, userID string) (*model.Collector, error)
	GetTrashItemByID(ctx context.Context, id string) (*model.AvaibleTrashByCollector, error)
	UpdateCollector(ctx context.Context, collector *model.Collector, updates map[string]interface{}) error
	UpdateAvaibleTrashByCollector(ctx context.Context, collectorID string, updatedTrash []model.AvaibleTrashByCollector) error
	DeleteAvaibleTrash(ctx context.Context, trashID string) error
}

type collectorRepository struct {
	// DB *gorm.DB
}

//	func NewCollectorRepository(db *gorm.DB) CollectorRepository {
//		return &collectorRepository{DB: db}
//	}
func NewCollectorRepository() CollectorRepository {
	return &collectorRepository{}
}

// func (r *collectorRepository) FindActiveCollectors() ([]model.Collector, error) {
// 	var collectors []model.Collector

// 	err := r.DB.Preload("Address").Where("job_status = ?", "active").First(&collectors).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch active collectors: %v", err)
// 	}

// 	return collectors, nil
// }

// func (r *collectorRepository) FindCollectorById(collector_id string) (*model.Collector, error) {
// 	var collector model.Collector
// 	err := r.DB.Preload("Address").Where("user_id = ?", collector_id).First(&collector).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching collector: %v", err)
// 	}
// 	fmt.Printf("menampilkan data collector %v", &collector)
// 	return &collector, nil
// }

// func (r *collectorRepository) FindCollectorByIdWithoutAddr(collector_id string) (*model.Collector, error) {
// 	var collector model.Collector
// 	err := r.DB.Where("user_id = ?", collector_id).First(&collector).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching collector: %v", err)
// 	}
// 	fmt.Printf("menampilkan data collector %v", &collector)
// 	return &collector, nil
// }

// func (r *collectorRepository) CreateCollector(collector *model.Collector) error {
// 	if err := r.DB.Create(collector).Error; err != nil {
// 		return fmt.Errorf("failed to create collector: %v", err)
// 	}
// 	return nil
// }

// func (r *collectorRepository) UpdateCollector(userId string, jobStatus string) (*model.Collector, error) {
// 	var existingCollector model.Collector

// 	if err := r.DB.Where("user_id = ?", userId).First(&existingCollector).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, fmt.Errorf("collector dengan user_id %s tidak ditemukan", userId)
// 		}
// 		log.Printf("Gagal mencari collector: %v", err)
// 		return nil, fmt.Errorf("gagal fetching collector: %w", err)
// 	}

// 	if jobStatus != "active" && jobStatus != "nonactive" {
// 		return nil, fmt.Errorf("invalid job status: %v", jobStatus)
// 	}

// 	if err := r.DB.Model(&existingCollector).Update("jobstatus", jobStatus).Error; err != nil {
// 		log.Printf("Gagal mengupdate data collector: %v", err)
// 		return nil, fmt.Errorf("gagal mengupdate job status untuk collector: %w", err)
// 	}

// 	return &existingCollector, nil
// }

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
