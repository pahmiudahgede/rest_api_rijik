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
