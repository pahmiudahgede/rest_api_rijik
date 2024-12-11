package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetBanners() ([]domain.Banner, error) {
	return repositories.GetBanners()
}

func GetBannerByID(id string) (domain.Banner, error) {
	banner, err := repositories.GetBannerByID(id)
	if err != nil {

		return domain.Banner{}, errors.New("banner not found")
	}
	return banner, nil
}

func CreateBanner(bannerName, bannerImage string) (domain.Banner, error) {
	newBanner := domain.Banner{
		BannerName:  bannerName,
		BannerImage: bannerImage,
	}

	if err := repositories.CreateBanner(&newBanner); err != nil {
		return domain.Banner{}, err
	}

	return newBanner, nil
}

func UpdateBanner(id, bannerName, bannerImage string) (domain.Banner, error) {

	banner, err := repositories.GetBannerByID(id)
	if err != nil {
		return domain.Banner{}, err
	}

	banner.BannerName = bannerName
	banner.BannerImage = bannerImage

	if err := repositories.UpdateBanner(&banner); err != nil {
		return domain.Banner{}, err
	}

	return banner, nil
}
