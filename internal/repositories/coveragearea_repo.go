package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type CoverageAreaRepository interface {
	FindCoverageByProvinceAndRegency(province, regency string) (*model.CoverageArea, error)
	CreateCoverage(coverage *model.CoverageArea) error
	FindCoverageById(id string) (*model.CoverageArea, error)
	FindAllCoverage() ([]model.CoverageArea, error)
	UpdateCoverage(id string, coverage *model.CoverageArea) error
	DeleteCoverage(id string) error
}

type coverageAreaRepository struct {
	DB *gorm.DB
}

func NewCoverageAreaRepository(db *gorm.DB) CoverageAreaRepository {
	return &coverageAreaRepository{DB: db}
}

func (r *coverageAreaRepository) FindCoverageByProvinceAndRegency(province, regency string) (*model.CoverageArea, error) {
	var coverage model.CoverageArea
	err := r.DB.Where("province = ? AND regency = ?", province, regency).First(&coverage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &coverage, nil
}

func (r *coverageAreaRepository) CreateCoverage(coverage *model.CoverageArea) error {
	if err := r.DB.Create(coverage).Error; err != nil {
		return fmt.Errorf("failed to create coverage: %v", err)
	}
	return nil
}

func (r *coverageAreaRepository) FindCoverageById(id string) (*model.CoverageArea, error) {
	var coverage model.CoverageArea
	err := r.DB.Where("id = ?", id).First(&coverage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("coverage with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch coverage by ID: %v", err)
	}
	return &coverage, nil
}

func (r *coverageAreaRepository) FindAllCoverage() ([]model.CoverageArea, error) {
	var coverage []model.CoverageArea
	err := r.DB.Find(&coverage).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coverage: %v", err)
	}

	return coverage, nil
}

func (r *coverageAreaRepository) UpdateCoverage(id string, coverage *model.CoverageArea) error {
	err := r.DB.Model(&model.CoverageArea{}).Where("id = ?", id).Updates(coverage).Error
	if err != nil {
		return fmt.Errorf("failed to update coverage: %v", err)
	}
	return nil
}

func (r *coverageAreaRepository) DeleteCoverage(id string) error {
	result := r.DB.Delete(&model.CoverageArea{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete coverage: %v", result.Error)
	}
	return nil
}
