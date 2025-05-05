package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type CompanyProfileRepository interface {
	CreateCompanyProfile(companyProfile *model.CompanyProfile) (*model.CompanyProfile, error)
	GetCompanyProfileByID(id string) (*model.CompanyProfile, error)
	GetCompanyProfilesByUserID(userID string) ([]model.CompanyProfile, error)
	UpdateCompanyProfile(id string, companyProfile *model.CompanyProfile) (*model.CompanyProfile, error)
	DeleteCompanyProfile(id string) error
}

type companyProfileRepository struct {
	DB *gorm.DB
}

func NewCompanyProfileRepository(db *gorm.DB) CompanyProfileRepository {
	return &companyProfileRepository{
		DB: db,
	}
}

func (r *companyProfileRepository) CreateCompanyProfile(companyProfile *model.CompanyProfile) (*model.CompanyProfile, error) {
	err := r.DB.Create(companyProfile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create company profile: %v", err)
	}
	return companyProfile, nil
}

func (r *companyProfileRepository) GetCompanyProfileByID(id string) (*model.CompanyProfile, error) {
	var companyProfile model.CompanyProfile
	err := r.DB.Preload("User").First(&companyProfile, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company profile with ID %s not found", id)
		}
		return nil, fmt.Errorf("error fetching company profile: %v", err)
	}
	return &companyProfile, nil
}

func (r *companyProfileRepository) GetCompanyProfilesByUserID(userID string) ([]model.CompanyProfile, error) {
	var companyProfiles []model.CompanyProfile
	err := r.DB.Preload("User").Where("user_id = ?", userID).Find(&companyProfiles).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching company profiles for userID %s: %v", userID, err)
	}
	return companyProfiles, nil
}

func (r *companyProfileRepository) UpdateCompanyProfile(id string, companyProfile *model.CompanyProfile) (*model.CompanyProfile, error) {
	var existingProfile model.CompanyProfile
	err := r.DB.First(&existingProfile, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("company profile not found: %v", err)
	}

	err = r.DB.Model(&existingProfile).Updates(companyProfile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update company profile: %v", err)
	}

	return &existingProfile, nil
}

func (r *companyProfileRepository) DeleteCompanyProfile(id string) error {
	err := r.DB.Delete(&model.CompanyProfile{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("failed to delete company profile: %v", err)
	}
	return nil
}
