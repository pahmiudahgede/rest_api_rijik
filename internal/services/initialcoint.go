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

type PointService struct {
	repo *repositories.PointRepository
}

func NewPointService(repo *repositories.PointRepository) *PointService {
	return &PointService{repo: repo}
}

func (s *PointService) GetAllPoints() ([]dto.PointResponse, error) {
	ctx := config.Context()

	cacheKey := "points:all"
	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedPoints []dto.PointResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedPoints); err == nil {
			return cachedPoints, nil
		}
	}

	points, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.PointResponse
	for _, point := range points {
		result = append(result, dto.PointResponse{
			ID:           point.ID,
			CoinName:     point.CoinName,
			ValuePerUnit: point.ValuePerUnit,
			CreatedAt:    utils.FormatDateToIndonesianFormat(point.CreatedAt),
			UpdatedAt:    utils.FormatDateToIndonesianFormat(point.UpdatedAt),
		})
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *PointService) GetPointByID(id string) (*dto.PointResponse, error) {
	ctx := config.Context()

	cacheKey := "points:" + id
	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedPoint dto.PointResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedPoint); err == nil {
			return &cachedPoint, nil
		}
	}

	point, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := &dto.PointResponse{
		ID:           point.ID,
		CoinName:     point.CoinName,
		ValuePerUnit: point.ValuePerUnit,
		CreatedAt:    utils.FormatDateToIndonesianFormat(point.CreatedAt),
		UpdatedAt:    utils.FormatDateToIndonesianFormat(point.UpdatedAt),
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *PointService) CreatePoint(request *dto.PointCreateRequest) (*dto.PointResponse, error) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	newPoint := &domain.Point{
		CoinName:     request.CoinName,
		ValuePerUnit: request.ValuePerUnit,
	}

	err := s.repo.Create(newPoint)
	if err != nil {
		return nil, err
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "points:all")

	response := &dto.PointResponse{
		ID:           newPoint.ID,
		CoinName:     newPoint.CoinName,
		ValuePerUnit: newPoint.ValuePerUnit,
		CreatedAt:    utils.FormatDateToIndonesianFormat(newPoint.CreatedAt),
		UpdatedAt:    utils.FormatDateToIndonesianFormat(newPoint.UpdatedAt),
	}

	return response, nil
}

func (s *PointService) UpdatePoint(id string, request *dto.PointUpdateRequest) (*dto.PointResponse, error) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	point, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("point not found")
	}

	point.CoinName = request.CoinName
	point.ValuePerUnit = request.ValuePerUnit
	point.UpdatedAt = time.Now()

	err = s.repo.Update(point)
	if err != nil {
		return nil, errors.New("failed to update point")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "points:all")
	config.RedisClient.Del(ctx, "points:"+id)

	response := &dto.PointResponse{
		ID:           point.ID,
		CoinName:     point.CoinName,
		ValuePerUnit: point.ValuePerUnit,
		CreatedAt:    utils.FormatDateToIndonesianFormat(point.CreatedAt),
		UpdatedAt:    utils.FormatDateToIndonesianFormat(point.UpdatedAt),
	}

	return response, nil
}

func (s *PointService) DeletePoint(id string) error {

	point, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("point not found")
	}

	err = s.repo.Delete(point)
	if err != nil {
		return errors.New("failed to delete point")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "points:all")
	config.RedisClient.Del(ctx, "points:"+id)

	return nil
}
