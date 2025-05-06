package repositories

import (
	"errors"
	"fmt"
	"log"

	"rijig/model"

	"gorm.io/gorm"
)

type TrashRepository interface {
	CreateCategory(category *model.TrashCategory) error
	AddDetailToCategory(detail *model.TrashDetail) error
	GetCategories() ([]model.TrashCategory, error)
	GetCategoryByID(id string) (*model.TrashCategory, error)
	GetTrashDetailByID(id string) (*model.TrashDetail, error)
	GetDetailsByCategoryID(categoryID string) ([]model.TrashDetail, error)
	UpdateCategoryName(id string, newName string) error
	UpdateCategory(id string, updateTrashCategory *model.TrashCategory) (*model.TrashCategory, error)
	UpdateTrashDetail(id string, description string, price float64) error
	DeleteCategory(id string) error
	DeleteTrashDetail(id string) error
}

type trashRepository struct {
	DB *gorm.DB
}

func NewTrashRepository(db *gorm.DB) TrashRepository {
	return &trashRepository{DB: db}
}

func (r *trashRepository) CreateCategory(category *model.TrashCategory) error {
	if err := r.DB.Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}
	return nil
}

func (r *trashRepository) AddDetailToCategory(detail *model.TrashDetail) error {
	if err := r.DB.Create(detail).Error; err != nil {
		return fmt.Errorf("failed to add detail to category: %v", err)
	}
	return nil
}

func (r *trashRepository) GetCategories() ([]model.TrashCategory, error) {
	var categories []model.TrashCategory
	if err := r.DB.Preload("Details").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %v", err)
	}
	return categories, nil
}

func (r *trashRepository) GetCategoryByID(id string) (*model.TrashCategory, error) {
	var category model.TrashCategory

	if err := r.DB.Preload("Details").First(&category, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}
	return &category, nil
}

func (r *trashRepository) GetTrashDetailByID(id string) (*model.TrashDetail, error) {
	var detail model.TrashDetail
	if err := r.DB.First(&detail, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("trash detail not found: %v", err)
	}
	return &detail, nil
}

func (r *trashRepository) GetDetailsByCategoryID(categoryID string) ([]model.TrashDetail, error) {
	var details []model.TrashDetail

	if err := r.DB.Where("category_id = ?", categoryID).Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch details for category %s: %v", categoryID, err)
	}
	return details, nil
}

func (r *trashRepository) UpdateCategoryName(id string, newName string) error {
	if err := r.DB.Model(&model.TrashCategory{}).Where("id = ?", id).Update("name", newName).Error; err != nil {
		return fmt.Errorf("failed to update category name: %v", err)
	}
	return nil
}

func (r *trashRepository) UpdateCategory(id string, updateTrashCategory *model.TrashCategory) (*model.TrashCategory, error) {
	var existingtrashCtgry model.TrashCategory
	if err := r.DB.Where("id = ?", id).First(&existingtrashCtgry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("trashCategory with ID %s not found", id)
		}
		log.Printf("Error fetching trash category for update: %v", err)
		return nil, fmt.Errorf("error fetching trash category for update: %w", err)
	}

	if err := r.DB.Save(&existingtrashCtgry).Error; err != nil {
		log.Printf("Error updating trash category: %v", err)
		return nil, fmt.Errorf("failed to update trash category: %w", err)
	}
	return &existingtrashCtgry, nil
}

func (r *trashRepository) UpdateTrashDetail(id string, description string, price float64) error {
	if err := r.DB.Model(&model.TrashDetail{}).Where("id = ?", id).Updates(model.TrashDetail{Description: description, Price: price}).Error; err != nil {
		return fmt.Errorf("failed to update trash detail: %v", err)
	}
	return nil
}

func (r *trashRepository) DeleteCategory(id string) error {
	if err := r.DB.Delete(&model.TrashCategory{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}
	return nil
}

func (r *trashRepository) DeleteTrashDetail(id string) error {
	if err := r.DB.Delete(&model.TrashDetail{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete trash detail: %v", err)
	}
	return nil
}
