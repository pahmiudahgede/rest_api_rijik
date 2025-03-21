package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type InitialCointRepository interface {
	CreateInitialCoint(coint *model.InitialCoint) error
	FindInitialCointByID(id string) (*model.InitialCoint, error)
	FindAllInitialCoints() ([]model.InitialCoint, error)
	UpdateInitialCoint(id string, coint *model.InitialCoint) error
	DeleteInitialCoint(id string) error
}

type initialCointRepository struct {
	DB *gorm.DB
}

func NewInitialCointRepository(db *gorm.DB) InitialCointRepository {
	return &initialCointRepository{DB: db}
}

func (r *initialCointRepository) CreateInitialCoint(coint *model.InitialCoint) error {
	if err := r.DB.Create(coint).Error; err != nil {
		return fmt.Errorf("failed to create initial coint: %v", err)
	}
	return nil
}

func (r *initialCointRepository) FindInitialCointByID(id string) (*model.InitialCoint, error) {
	var coint model.InitialCoint
	err := r.DB.Where("id = ?", id).First(&coint).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("initial coint with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch initial coint by ID: %v", err)
	}
	return &coint, nil
}

func (r *initialCointRepository) FindAllInitialCoints() ([]model.InitialCoint, error) {
	var coints []model.InitialCoint
	err := r.DB.Find(&coints).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial coints: %v", err)
	}
	return coints, nil
}

func (r *initialCointRepository) UpdateInitialCoint(id string, coint *model.InitialCoint) error {
	err := r.DB.Model(&model.InitialCoint{}).Where("id = ?", id).Updates(coint).Error
	if err != nil {
		return fmt.Errorf("failed to update initial coint: %v", err)
	}
	return nil
}

func (r *initialCointRepository) DeleteInitialCoint(id string) error {
	result := r.DB.Delete(&model.InitialCoint{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete initial coint: %v", result.Error)
	}
	return nil
}
