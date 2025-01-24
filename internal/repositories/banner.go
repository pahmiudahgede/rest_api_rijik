package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

type BannerRepository struct{}

func NewBannerRepository() *BannerRepository {
	return &BannerRepository{}
}

func (r *BannerRepository) GetAll() ([]domain.Banner, error) {
	var banners []domain.Banner
	err := config.DB.Find(&banners).Error
	if err != nil {
		return nil, errors.New("failed to fetch banners from database")
	}
	return banners, nil
}

func (r *BannerRepository) GetByID(id string) (*domain.Banner, error) {
	var banner domain.Banner
	err := config.DB.First(&banner, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("banner not found")
	}
	return &banner, nil
}

func (r *BannerRepository) Create(banner *domain.Banner) error {
	return config.DB.Create(banner).Error
}

func (r *BannerRepository) Update(banner *domain.Banner) error {
	return config.DB.Save(banner).Error
}

func (r *BannerRepository) Delete(banner *domain.Banner) error {
	return config.DB.Delete(banner).Error
}
