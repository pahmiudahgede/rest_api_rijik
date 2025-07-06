package trash

import (
	"context"
	"errors"
	"fmt"
	"rijig/model"
	"time"

	"gorm.io/gorm"
)

type TrashRepositoryInterface interface {
	CreateTrashCategory(ctx context.Context, category *model.TrashCategory) error
	CreateTrashCategoryWithDetails(ctx context.Context, category *model.TrashCategory, details []model.TrashDetail) error
	UpdateTrashCategory(ctx context.Context, id string, updates map[string]interface{}) error
	GetAllTrashCategories(ctx context.Context) ([]model.TrashCategory, error)
	GetAllTrashCategoriesWithDetails(ctx context.Context) ([]model.TrashCategory, error)
	GetTrashCategoryByID(ctx context.Context, id string) (*model.TrashCategory, error)
	GetTrashCategoryByIDWithDetails(ctx context.Context, id string) (*model.TrashCategory, error)
	DeleteTrashCategory(ctx context.Context, id string) error

	CreateTrashDetail(ctx context.Context, detail *model.TrashDetail) error
	AddTrashDetailToCategory(ctx context.Context, categoryID string, detail *model.TrashDetail) error
	UpdateTrashDetail(ctx context.Context, id string, updates map[string]interface{}) error
	GetTrashDetailsByCategory(ctx context.Context, categoryID string) ([]model.TrashDetail, error)
	GetTrashDetailByID(ctx context.Context, id string) (*model.TrashDetail, error)
	DeleteTrashDetail(ctx context.Context, id string) error

	CheckTrashCategoryExists(ctx context.Context, id string) (bool, error)
	CheckTrashDetailExists(ctx context.Context, id string) (bool, error)
	GetMaxStepOrderByCategory(ctx context.Context, categoryID string) (int, error)
}

type trashRepository struct {
	db *gorm.DB
}

func NewTrashRepository(db *gorm.DB) TrashRepositoryInterface {
	return &trashRepository{
		db,
	}
}

func (r *trashRepository) CreateTrashCategory(ctx context.Context, category *model.TrashCategory) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create trash category: %w", err)
	}
	return nil
}

func (r *trashRepository) CreateTrashCategoryWithDetails(ctx context.Context, category *model.TrashCategory, details []model.TrashDetail) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(category).Error; err != nil {
			return fmt.Errorf("failed to create trash category: %w", err)
		}

		if len(details) > 0 {

			for i := range details {
				details[i].TrashCategoryID = category.ID

				if details[i].StepOrder == 0 {
					details[i].StepOrder = i + 1
				}
			}

			if err := tx.Create(&details).Error; err != nil {
				return fmt.Errorf("failed to create trash details: %w", err)
			}
		}

		return nil
	})
}

func (r *trashRepository) UpdateTrashCategory(ctx context.Context, id string, updates map[string]interface{}) error {

	exists, err := r.CheckTrashCategoryExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash category not found")
	}

	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&model.TrashCategory{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update trash category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected during update")
	}

	return nil
}

func (r *trashRepository) GetAllTrashCategories(ctx context.Context) ([]model.TrashCategory, error) {
	var categories []model.TrashCategory

	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get trash categories: %w", err)
	}

	return categories, nil
}

func (r *trashRepository) GetAllTrashCategoriesWithDetails(ctx context.Context) ([]model.TrashCategory, error) {
	var categories []model.TrashCategory

	if err := r.db.WithContext(ctx).Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get trash categories with details: %w", err)
	}

	return categories, nil
}

func (r *trashRepository) GetTrashCategoryByID(ctx context.Context, id string) (*model.TrashCategory, error) {
	var category model.TrashCategory

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trash category not found")
		}
		return nil, fmt.Errorf("failed to get trash category: %w", err)
	}

	return &category, nil
}

func (r *trashRepository) GetTrashCategoryByIDWithDetails(ctx context.Context, id string) (*model.TrashCategory, error) {
	var category model.TrashCategory

	if err := r.db.WithContext(ctx).Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trash category not found")
		}
		return nil, fmt.Errorf("failed to get trash category with details: %w", err)
	}

	return &category, nil
}

func (r *trashRepository) DeleteTrashCategory(ctx context.Context, id string) error {

	exists, err := r.CheckTrashCategoryExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash category not found")
	}

	result := r.db.WithContext(ctx).Delete(&model.TrashCategory{ID: id})
	if result.Error != nil {
		return fmt.Errorf("failed to delete trash category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected during deletion")
	}

	return nil
}

func (r *trashRepository) CreateTrashDetail(ctx context.Context, detail *model.TrashDetail) error {

	exists, err := r.CheckTrashCategoryExists(ctx, detail.TrashCategoryID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash category not found")
	}

	if detail.StepOrder == 0 {
		maxOrder, err := r.GetMaxStepOrderByCategory(ctx, detail.TrashCategoryID)
		if err != nil {
			return err
		}
		detail.StepOrder = maxOrder + 1
	}

	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return fmt.Errorf("failed to create trash detail: %w", err)
	}

	return nil
}

func (r *trashRepository) AddTrashDetailToCategory(ctx context.Context, categoryID string, detail *model.TrashDetail) error {

	exists, err := r.CheckTrashCategoryExists(ctx, categoryID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash category not found")
	}

	detail.TrashCategoryID = categoryID

	if detail.StepOrder == 0 {
		maxOrder, err := r.GetMaxStepOrderByCategory(ctx, categoryID)
		if err != nil {
			return err
		}
		detail.StepOrder = maxOrder + 1
	}

	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return fmt.Errorf("failed to add trash detail to category: %w", err)
	}

	return nil
}

func (r *trashRepository) UpdateTrashDetail(ctx context.Context, id string, updates map[string]interface{}) error {

	exists, err := r.CheckTrashDetailExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash detail not found")
	}

	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&model.TrashDetail{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update trash detail: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected during update")
	}

	return nil
}

func (r *trashRepository) GetTrashDetailsByCategory(ctx context.Context, categoryID string) ([]model.TrashDetail, error) {
	var details []model.TrashDetail

	if err := r.db.WithContext(ctx).Where("trash_category_id = ?", categoryID).Order("step_order ASC").Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to get trash details: %w", err)
	}

	return details, nil
}

func (r *trashRepository) GetTrashDetailByID(ctx context.Context, id string) (*model.TrashDetail, error) {
	var detail model.TrashDetail

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trash detail not found")
		}
		return nil, fmt.Errorf("failed to get trash detail: %w", err)
	}

	return &detail, nil
}

func (r *trashRepository) DeleteTrashDetail(ctx context.Context, id string) error {

	exists, err := r.CheckTrashDetailExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trash detail not found")
	}

	result := r.db.WithContext(ctx).Delete(&model.TrashDetail{ID: id})
	if result.Error != nil {
		return fmt.Errorf("failed to delete trash detail: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected during deletion")
	}

	return nil
}

func (r *trashRepository) CheckTrashCategoryExists(ctx context.Context, id string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.TrashCategory{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check trash category existence: %w", err)
	}

	return count > 0, nil
}

func (r *trashRepository) CheckTrashDetailExists(ctx context.Context, id string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.TrashDetail{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check trash detail existence: %w", err)
	}

	return count > 0, nil
}

func (r *trashRepository) GetMaxStepOrderByCategory(ctx context.Context, categoryID string) (int, error) {
	var maxOrder int

	if err := r.db.WithContext(ctx).Model(&model.TrashDetail{}).
		Where("trash_category_id = ?", categoryID).
		Select("COALESCE(MAX(step_order), 0)").
		Scan(&maxOrder).Error; err != nil {
		return 0, fmt.Errorf("failed to get max step order: %w", err)
	}

	return maxOrder, nil
}
