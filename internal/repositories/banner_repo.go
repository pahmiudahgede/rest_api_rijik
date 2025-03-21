package repositories

import (
	"fmt"

	"rijig/model"

	"gorm.io/gorm"
)

type BannerRepository interface {
	CreateBanner(banner *model.Banner) error
	FindBannerByID(id string) (*model.Banner, error)
	FindAllBanners() ([]model.Banner, error)
	UpdateBanner(id string, banner *model.Banner) error
	DeleteBanner(id string) error
}

type bannerRepository struct {
	DB *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{DB: db}
}

func (r *bannerRepository) CreateBanner(banner *model.Banner) error {
	if err := r.DB.Create(banner).Error; err != nil {
		return fmt.Errorf("failed to create banner: %v", err)
	}
	return nil
}

func (r *bannerRepository) FindBannerByID(id string) (*model.Banner, error) {
	var banner model.Banner
	err := r.DB.Where("id = ?", id).First(&banner).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("banner with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch banner by ID: %v", err)
	}
	return &banner, nil
}

func (r *bannerRepository) FindAllBanners() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.DB.Find(&banners).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banners: %v", err)
	}

	return banners, nil
}

func (r *bannerRepository) UpdateBanner(id string, banner *model.Banner) error {
	err := r.DB.Model(&model.Banner{}).Where("id = ?", id).Updates(banner).Error
	if err != nil {
		return fmt.Errorf("failed to update banner: %v", err)
	}
	return nil
}

func (r *bannerRepository) DeleteBanner(id string) error {
	result := r.DB.Delete(&model.Banner{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete banner: %v", result.Error)
	}
	return nil
}
