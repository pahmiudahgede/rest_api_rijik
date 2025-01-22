package repositories

import (
	"encoding/json"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetBanners() ([]domain.Banner, error) {
	var banners []domain.Banner

	ctx := config.Context()
	cachedBanners, err := config.RedisClient.Get(ctx, "banners").Result()
	if err == nil && cachedBanners != "" {

		if err := json.Unmarshal([]byte(cachedBanners), &banners); err != nil {
			return nil, err
		}
		return banners, nil
	}

	if err := config.DB.Find(&banners).Error; err != nil {
		return nil, err
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		return nil, err
	}
	config.RedisClient.Set(ctx, "banners", bannersJSON, time.Hour).Err()

	return banners, nil
}

func GetBannerByID(id string) (domain.Banner, error) {
	var banner domain.Banner

	ctx := config.Context()
	cachedBanner, err := config.RedisClient.Get(ctx, "banner:"+id).Result()
	if err == nil && cachedBanner != "" {

		if err := json.Unmarshal([]byte(cachedBanner), &banner); err != nil {
			return banner, err
		}
		return banner, nil
	}

	if err := config.DB.Where("id = ?", id).First(&banner).Error; err != nil {
		return banner, err
	}

	bannerJSON, err := json.Marshal(banner)
	if err != nil {
		return banner, err
	}
	config.RedisClient.Set(ctx, "banner:"+id, bannerJSON, time.Hour*24).Err()

	return banner, nil
}

func CreateBanner(banner *domain.Banner) error {
	if err := config.DB.Create(banner).Error; err != nil {
		return err
	}

	config.RedisClient.Del(config.Context(), "banners")
	return nil
}

func UpdateBanner(banner *domain.Banner) error {
	if err := config.DB.Save(banner).Error; err != nil {
		return err
	}

	config.RedisClient.Del(config.Context(), "banners")
	config.RedisClient.Del(config.Context(), "banner:"+banner.ID)
	return nil
}

func DeleteBanner(id string) error {
	if err := config.DB.Where("id = ?", id).Delete(&domain.Banner{}).Error; err != nil {
		return err
	}

	config.RedisClient.Del(config.Context(), "banners")
	config.RedisClient.Del(config.Context(), "banner:"+id)
	return nil
}
