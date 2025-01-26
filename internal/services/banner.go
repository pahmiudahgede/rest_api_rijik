package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type BannerService struct {
	repo *repositories.BannerRepository
}

func NewBannerService(repo *repositories.BannerRepository) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) GetAllBanners() ([]dto.BannerResponse, error) {
	ctx := config.Context()
	cacheKey := "banners:all"

	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedBanners []dto.BannerResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedBanners); err == nil {
			return cachedBanners, nil
		}
	}

	banners, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.BannerResponse
	for _, banner := range banners {
		result = append(result, dto.BannerResponse{
			ID:          banner.ID,
			BannerName:  banner.BannerName,
			BannerImage: banner.BannerImage,
			CreatedAt:   utils.FormatDateToIndonesianFormat(banner.CreatedAt),
			UpdatedAt:   utils.FormatDateToIndonesianFormat(banner.UpdatedAt),
		})
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *BannerService) GetBannerByID(id string) (*dto.BannerResponse, error) {
	ctx := config.Context()
	cacheKey := "banners:" + id

	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedBanner dto.BannerResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedBanner); err == nil {
			return &cachedBanner, nil
		}
	}

	banner, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := &dto.BannerResponse{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   utils.FormatDateToIndonesianFormat(banner.CreatedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(banner.UpdatedAt),
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *BannerService) CreateBanner(request *dto.BannerCreateRequest) (*dto.BannerResponse, error) {

	if request.BannerName == "" || request.BannerImage == "" {
		return nil, errors.New("invalid input data")
	}

	newBanner := &domain.Banner{
		BannerName:  request.BannerName,
		BannerImage: request.BannerImage,
	}

	err := s.repo.Create(newBanner)
	if err != nil {
		return nil, errors.New("failed to create banner")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "banners:all")

	response := &dto.BannerResponse{
		ID:          newBanner.ID,
		BannerName:  newBanner.BannerName,
		BannerImage: newBanner.BannerImage,
		CreatedAt:   utils.FormatDateToIndonesianFormat(newBanner.CreatedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(newBanner.UpdatedAt),
	}

	return response, nil
}

func (s *BannerService) UpdateBanner(id string, request *dto.BannerUpdateRequest) (*dto.BannerResponse, error) {

	banner, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("banner not found")
	}

	if request.BannerName != nil && *request.BannerName != "" {
		banner.BannerName = *request.BannerName
	}
	if request.BannerImage != nil && *request.BannerImage != "" {
		banner.BannerImage = *request.BannerImage
	}
	banner.UpdatedAt = time.Now()

	err = s.repo.Update(banner)
	if err != nil {
		return nil, errors.New("failed to update banner")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "banners:all")
	config.RedisClient.Del(ctx, "banners:"+id)

	response := &dto.BannerResponse{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   utils.FormatDateToIndonesianFormat(banner.CreatedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(banner.UpdatedAt),
	}

	return response, nil
}

func (s *BannerService) DeleteBanner(id string) error {

	banner, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("banner not found")
	}

	err = s.repo.Delete(banner)
	if err != nil {
		return errors.New("failed to delete banner")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "banners:all")
	config.RedisClient.Del(ctx, "banners:"+id)

	return nil
}
