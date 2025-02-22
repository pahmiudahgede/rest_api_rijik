package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type InitialCointService interface {
	CreateInitialCoint(request dto.RequestInitialCointDTO) (*dto.ReponseInitialCointDTO, error)
	GetAllInitialCoints() ([]dto.ReponseInitialCointDTO, error)
	GetInitialCointByID(id string) (*dto.ReponseInitialCointDTO, error)
	UpdateInitialCoint(id string, request dto.RequestInitialCointDTO) (*dto.ReponseInitialCointDTO, error)
	DeleteInitialCoint(id string) error
}

type initialCointService struct {
	InitialCointRepo repositories.InitialCointRepository
}

func NewInitialCointService(repo repositories.InitialCointRepository) InitialCointService {
	return &initialCointService{InitialCointRepo: repo}
}

func (s *initialCointService) CreateInitialCoint(request dto.RequestInitialCointDTO) (*dto.ReponseInitialCointDTO, error) {

	errors, valid := request.ValidateCointInput()
	if !valid {
		return nil, fmt.Errorf("validation error: %v", errors)
	}

	coint := model.InitialCoint{
		CoinName:     request.CoinName,
		ValuePerUnit: request.ValuePerUnit,
	}
	if err := s.InitialCointRepo.CreateInitialCoint(&coint); err != nil {
		return nil, fmt.Errorf("failed to create initial coint: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(coint.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(coint.UpdatedAt)

	responseDTO := &dto.ReponseInitialCointDTO{
		ID:           coint.ID,
		CoinName:     coint.CoinName,
		ValuePerUnit: coint.ValuePerUnit,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	cacheKey := fmt.Sprintf("initialcoint:%s", coint.ID)
	cacheData := map[string]interface{}{
		"data": responseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching new initial coint: %v\n", err)
	}

	err := s.updateAllCointCache()
	if err != nil {
		return nil, fmt.Errorf("error updating all initial coint cache: %v", err)
	}

	return responseDTO, nil
}

func (s *initialCointService) GetAllInitialCoints() ([]dto.ReponseInitialCointDTO, error) {
	var cointsDTO []dto.ReponseInitialCointDTO
	cacheKey := "initialcoints:all"

	cachedData, err := utils.GetJSONData(cacheKey)
	if err != nil {
		fmt.Printf("Error fetching cache for initialcoints: %v\n", err)
	}

	if cachedData != nil {
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				if cointData, ok := item.(map[string]interface{}); ok {

					if coinID, ok := cointData["coin_id"].(string); ok {
						if coinName, ok := cointData["coin_name"].(string); ok {
							if valuePerUnit, ok := cointData["value_perunit"].(float64); ok {
								if createdAt, ok := cointData["createdAt"].(string); ok {
									if updatedAt, ok := cointData["updatedAt"].(string); ok {

										cointsDTO = append(cointsDTO, dto.ReponseInitialCointDTO{
											ID:           coinID,
											CoinName:     coinName,
											ValuePerUnit: valuePerUnit,
											CreatedAt:    createdAt,
											UpdatedAt:    updatedAt,
										})
									}
								}
							}
						}
					}
				}
			}
			return cointsDTO, nil
		}
	}

	records, err := s.InitialCointRepo.FindAllInitialCoints()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial coints from database: %v", err)
	}

	if len(records) == 0 {
		return cointsDTO, nil
	}

	for _, record := range records {
		createdAt, _ := utils.FormatDateToIndonesianFormat(record.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(record.UpdatedAt)

		cointsDTO = append(cointsDTO, dto.ReponseInitialCointDTO{
			ID:           record.ID,
			CoinName:     record.CoinName,
			ValuePerUnit: record.ValuePerUnit,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": cointsDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching all initial coints: %v\n", err)
	}

	return cointsDTO, nil
}

func (s *initialCointService) GetInitialCointByID(id string) (*dto.ReponseInitialCointDTO, error) {
	cacheKey := fmt.Sprintf("initialcoint:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		if data, ok := cachedData["data"].(map[string]interface{}); ok {

			return &dto.ReponseInitialCointDTO{
				ID:           data["coin_id"].(string),
				CoinName:     data["coin_name"].(string),
				ValuePerUnit: data["value_perunit"].(float64),
				CreatedAt:    data["createdAt"].(string),
				UpdatedAt:    data["updatedAt"].(string),
			}, nil
		} else {
			return nil, fmt.Errorf("error: cache data is not in the expected format for coin ID %s", id)
		}
	}

	coint, err := s.InitialCointRepo.FindInitialCointByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial coint by ID %s: %v", id, err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(coint.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(coint.UpdatedAt)

	cointDTO := &dto.ReponseInitialCointDTO{
		ID:           coint.ID,
		CoinName:     coint.CoinName,
		ValuePerUnit: coint.ValuePerUnit,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": cointDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching initial coint by ID: %v\n", err)
	}

	return cointDTO, nil
}

func (s *initialCointService) UpdateInitialCoint(id string, request dto.RequestInitialCointDTO) (*dto.ReponseInitialCointDTO, error) {

	coint, err := s.InitialCointRepo.FindInitialCointByID(id)
	if err != nil {
		return nil, fmt.Errorf("initial coint with ID %s not found", id)
	}

	coint.CoinName = request.CoinName
	coint.ValuePerUnit = request.ValuePerUnit

	if err := s.InitialCointRepo.UpdateInitialCoint(id, coint); err != nil {
		return nil, fmt.Errorf("failed to update initial coint: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(coint.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(coint.UpdatedAt)

	cointDTO := &dto.ReponseInitialCointDTO{
		ID:           coint.ID,
		CoinName:     coint.CoinName,
		ValuePerUnit: coint.ValuePerUnit,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	cacheKey := fmt.Sprintf("initialcoint:%s", id)
	cacheData := map[string]interface{}{
		"data": cointDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching updated initial coint: %v\n", err)
	}

	allCoints, err := s.InitialCointRepo.FindAllInitialCoints()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all initial coints from database: %v", err)
	}

	var cointsDTO []dto.ReponseInitialCointDTO
	for _, record := range allCoints {
		createdAt, _ := utils.FormatDateToIndonesianFormat(record.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(record.UpdatedAt)

		cointsDTO = append(cointsDTO, dto.ReponseInitialCointDTO{
			ID:           record.ID,
			CoinName:     record.CoinName,
			ValuePerUnit: record.ValuePerUnit,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	cacheAllKey := "initialcoints:all"
	cacheAllData := map[string]interface{}{
		"data": cointsDTO,
	}
	if err := utils.SetJSONData(cacheAllKey, cacheAllData, time.Hour*24); err != nil {
		fmt.Printf("Error caching all initial coints: %v\n", err)
	}

	return cointDTO, nil
}

func (s *initialCointService) DeleteInitialCoint(id string) error {

	coint, err := s.InitialCointRepo.FindInitialCointByID(id)
	if err != nil {
		return fmt.Errorf("initial coint with ID %s not found", id)
	}

	if err := s.InitialCointRepo.DeleteInitialCoint(id); err != nil {
		return fmt.Errorf("failed to delete initial coint: %v", err)
	}

	cacheKey := fmt.Sprintf("initialcoint:%s", coint.ID)
	if err := utils.DeleteData(cacheKey); err != nil {
		fmt.Printf("Error deleting cache for initial coint: %v\n", err)
	}

	return s.updateAllCointCache()
}

func (s *initialCointService) updateAllCointCache() error {

	records, err := s.InitialCointRepo.FindAllInitialCoints()
	if err != nil {
		return fmt.Errorf("failed to fetch all initial coints from database: %v", err)
	}

	var cointsDTO []dto.ReponseInitialCointDTO
	for _, record := range records {
		createdAt, _ := utils.FormatDateToIndonesianFormat(record.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(record.UpdatedAt)

		cointsDTO = append(cointsDTO, dto.ReponseInitialCointDTO{
			ID:           record.ID,
			CoinName:     record.CoinName,
			ValuePerUnit: record.ValuePerUnit,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	cacheAllKey := "initialcoints:all"
	cacheAllData := map[string]interface{}{
		"data": cointsDTO,
	}
	if err := utils.SetJSONData(cacheAllKey, cacheAllData, time.Hour*24); err != nil {
		fmt.Printf("Error caching all initial coints: %v\n", err)
		return err
	}

	return nil
}
