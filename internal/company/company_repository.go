package company

import (
	"context"
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type CompanyProfileRepository interface {
	CreateCompanyProfile(ctx context.Context, companyProfile *model.CompanyProfile) (*model.CompanyProfile, error)
	GetCompanyProfileByID(ctx context.Context, id string) (*model.CompanyProfile, error)
	GetCompanyProfilesByUserID(ctx context.Context, userID string) ([]model.CompanyProfile, error)
	UpdateCompanyProfile(ctx context.Context, company *model.CompanyProfile) error
	DeleteCompanyProfileByUserID(ctx context.Context, userID string) error
	ExistsByUserID(ctx context.Context, userID string) (bool, error) 
}

type companyProfileRepository struct {
	db *gorm.DB
}

func NewCompanyProfileRepository(db *gorm.DB) CompanyProfileRepository {
	return &companyProfileRepository{db}
}

func (r *companyProfileRepository) CreateCompanyProfile(ctx context.Context, companyProfile *model.CompanyProfile) (*model.CompanyProfile, error) {
	err := r.db.WithContext(ctx).Create(companyProfile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create company profile: %v", err)
	}
	return companyProfile, nil
}

func (r *companyProfileRepository) GetCompanyProfileByID(ctx context.Context, id string) (*model.CompanyProfile, error) {
	var companyProfile model.CompanyProfile
	err := r.db.WithContext(ctx).Preload("User").First(&companyProfile, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company profile with ID %s not found", id)
		}
		return nil, fmt.Errorf("error fetching company profile: %v", err)
	}
	return &companyProfile, nil
}

func (r *companyProfileRepository) GetCompanyProfilesByUserID(ctx context.Context, userID string) ([]model.CompanyProfile, error) {
	var companyProfiles []model.CompanyProfile
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).Find(&companyProfiles).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching company profiles for userID %s: %v", userID, err)
	}
	return companyProfiles, nil
}

func (r *companyProfileRepository) UpdateCompanyProfile(ctx context.Context, company *model.CompanyProfile) error {
	var existing model.CompanyProfile
	if err := r.db.WithContext(ctx).First(&existing, "user_id = ?", company.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("company profile not found for user_id %s", company.UserID)
		}
		return fmt.Errorf("failed to fetch company profile: %v", err)
	}

	err := r.db.WithContext(ctx).Model(&existing).Updates(company).Error
	if err != nil {
		return fmt.Errorf("failed to update company profile: %v", err)
	}
	return nil
}

func (r *companyProfileRepository) DeleteCompanyProfileByUserID(ctx context.Context, userID string) error {
	err := r.db.WithContext(ctx).Delete(&model.CompanyProfile{}, "user_id = ?", userID).Error
	if err != nil {
		return fmt.Errorf("failed to delete company profile: %v", err)
	}
	return nil
}

func (r *companyProfileRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CompanyProfile{}).
		Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %v", err)
	}
	return count > 0, nil
}
