package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetBanners() ([]domain.Banner, error) {
	var banners []domain.Banner

	if err := config.DB.Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

func GetBannerByID(id string) (domain.Banner, error) {
	var banner domain.Banner

	if err := config.DB.Where("id = ?", id).First(&banner).Error; err != nil {
		return banner, err
	}
	return banner, nil
}

func CreateBanner(banner *domain.Banner) error {
	if err := config.DB.Create(banner).Error; err != nil {
		return err
	}
	return nil
}

func UpdateBanner(banner *domain.Banner) error {
	if err := config.DB.Save(banner).Error; err != nil {
		return err
	}
	return nil
}

func DeleteBanner(id string) error {
	if err := config.DB.Where("id = ?", id).Delete(&domain.Banner{}).Error; err != nil {
		return err
	}
	return nil
}