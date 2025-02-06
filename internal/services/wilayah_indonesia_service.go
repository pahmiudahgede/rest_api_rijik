package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type WilayahIndonesiaService interface {
	ImportDataFromCSV() error

	GetAllProvinces(page, limit int) ([]dto.ProvinceResponseDTO, error)
	GetProvinceByID(id string) (*dto.ProvinceResponseDTO, error)

	GetAllRegencies(page, limit int) ([]dto.RegencyResponseDTO, error)
	GetRegencyByID(id string) (*dto.RegencyResponseDTO, error)

	GetAllDistricts(page, limit int) ([]dto.DistrictResponseDTO, error)
	GetDistrictByID(id string) (*dto.DistrictResponseDTO, error)

	GetAllVillages(page, limit int) ([]dto.VillageResponseDTO, error)
	GetVillageByID(id string) (*dto.VillageResponseDTO, error)
}

type wilayahIndonesiaService struct {
	WilayahRepo repositories.WilayahIndonesiaRepository
}

func NewWilayahIndonesiaService(wilayahRepo repositories.WilayahIndonesiaRepository) WilayahIndonesiaService {
	return &wilayahIndonesiaService{WilayahRepo: wilayahRepo}
}

func (s *wilayahIndonesiaService) ImportDataFromCSV() error {

	provinces, err := utils.ReadCSV("public/document/provinces.csv")
	if err != nil {
		return fmt.Errorf("failed to read provinces CSV: %v", err)
	}

	var provinceList []model.Province
	for _, record := range provinces[1:] {
		province := model.Province{
			ID:   record[0],
			Name: record[1],
		}
		provinceList = append(provinceList, province)
	}

	if err := s.WilayahRepo.ImportProvinces(provinceList); err != nil {
		return fmt.Errorf("failed to import provinces: %v", err)
	}

	regencies, err := utils.ReadCSV("public/document/regencies.csv")
	if err != nil {
		return fmt.Errorf("failed to read regencies CSV: %v", err)
	}

	var regencyList []model.Regency
	for _, record := range regencies[1:] {
		regency := model.Regency{
			ID:         record[0],
			ProvinceID: record[1],
			Name:       record[2],
		}
		regencyList = append(regencyList, regency)
	}

	if err := s.WilayahRepo.ImportRegencies(regencyList); err != nil {
		return fmt.Errorf("failed to import regencies: %v", err)
	}

	districts, err := utils.ReadCSV("public/document/districts.csv")
	if err != nil {
		return fmt.Errorf("failed to read districts CSV: %v", err)
	}

	var districtList []model.District
	for _, record := range districts[1:] {
		district := model.District{
			ID:        record[0],
			RegencyID: record[1],
			Name:      record[2],
		}
		districtList = append(districtList, district)
	}

	if err := s.WilayahRepo.ImportDistricts(districtList); err != nil {
		return fmt.Errorf("failed to import districts: %v", err)
	}

	villages, err := utils.ReadCSV("public/document/villages.csv")
	if err != nil {
		return fmt.Errorf("failed to read villages CSV: %v", err)
	}

	var villageList []model.Village
	for _, record := range villages[1:] {
		village := model.Village{
			ID:         record[0],
			DistrictID: record[1],
			Name:       record[2],
		}
		villageList = append(villageList, village)
	}

	if err := s.WilayahRepo.ImportVillages(villageList); err != nil {
		return fmt.Errorf("failed to import villages: %v", err)
	}

	return nil
}

func (s *wilayahIndonesiaService) GetAllProvinces(page, limit int) ([]dto.ProvinceResponseDTO, error) {

	cacheKey := fmt.Sprintf("provinces_page_%d_limit_%d", page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var provinces []dto.ProvinceResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				province, ok := item.(map[string]interface{})
				if ok {
					provinces = append(provinces, dto.ProvinceResponseDTO{
						ID:   province["id"].(string),
						Name: province["name"].(string),
					})
				}
			}
			return provinces, nil
		}
	}

	provinces, total, err := s.WilayahRepo.FindAllProvinces(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %v", err)
	}

	var provinceDTOs []dto.ProvinceResponseDTO
	for _, province := range provinces {
		provinceDTOs = append(provinceDTOs, dto.ProvinceResponseDTO{
			ID:   province.ID,
			Name: province.Name,
		})
	}

	cacheData := map[string]interface{}{
		"data":  provinceDTOs,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching provinces data to Redis: %v\n", err)
	}

	return provinceDTOs, nil
}

func (s *wilayahIndonesiaService) GetProvinceByID(id string) (*dto.ProvinceResponseDTO, error) {

	cacheKey := fmt.Sprintf("province:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var province dto.ProvinceResponseDTO
		if data, ok := cachedData["data"].(map[string]interface{}); ok {
			province = dto.ProvinceResponseDTO{
				ID:   data["id"].(string),
				Name: data["name"].(string),
			}
			return &province, nil
		}
	}

	province, err := s.WilayahRepo.FindProvinceByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch province: %v", err)
	}

	provinceDTO := &dto.ProvinceResponseDTO{
		ID:   province.ID,
		Name: province.Name,
	}

	regenciesDTO := []dto.RegencyResponseDTO{}
	for _, regency := range province.Regencies {
		regenciesDTO = append(regenciesDTO, dto.RegencyResponseDTO{
			ID:         regency.ID,
			ProvinceID: regency.ProvinceID,
			Name:       regency.Name,
		})
	}
	provinceDTO.Regencies = regenciesDTO

	cacheData := map[string]interface{}{
		"data": provinceDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching province data to Redis: %v\n", err)
	}

	return provinceDTO, nil
}

func (s *wilayahIndonesiaService) GetAllRegencies(page, limit int) ([]dto.RegencyResponseDTO, error) {

	cacheKey := fmt.Sprintf("regencies_page_%d_limit_%d", page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var regencies []dto.RegencyResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				regency, ok := item.(map[string]interface{})
				if ok {
					regencies = append(regencies, dto.RegencyResponseDTO{
						ID:         regency["id"].(string),
						ProvinceID: regency["province_id"].(string),
						Name:       regency["name"].(string),
					})
				}
			}
			return regencies, nil
		}
	}

	regencies, total, err := s.WilayahRepo.FindAllRegencies(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch regencies: %v", err)
	}

	var regencyDTOs []dto.RegencyResponseDTO
	for _, regency := range regencies {
		regencyDTOs = append(regencyDTOs, dto.RegencyResponseDTO{
			ID:         regency.ID,
			ProvinceID: regency.ProvinceID,
			Name:       regency.Name,
		})
	}

	cacheData := map[string]interface{}{
		"data":  regencyDTOs,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching regencies data to Redis: %v\n", err)
	}

	return regencyDTOs, nil
}

func (s *wilayahIndonesiaService) GetRegencyByID(id string) (*dto.RegencyResponseDTO, error) {

	cacheKey := fmt.Sprintf("regency:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var regency dto.RegencyResponseDTO
		if data, ok := cachedData["data"].(map[string]interface{}); ok {
			regency = dto.RegencyResponseDTO{
				ID:         data["id"].(string),
				ProvinceID: data["province_id"].(string),
				Name:       data["name"].(string),
			}
			return &regency, nil
		}
	}

	regency, err := s.WilayahRepo.FindRegencyByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch regency: %v", err)
	}

	regencyDTO := &dto.RegencyResponseDTO{
		ID:         regency.ID,
		ProvinceID: regency.ProvinceID,
		Name:       regency.Name,
	}

	districtsDTO := []dto.DistrictResponseDTO{}
	for _, district := range regency.Districts {
		districtsDTO = append(districtsDTO, dto.DistrictResponseDTO{
			ID:        district.ID,
			RegencyID: district.RegencyID,
			Name:      district.Name,
		})
	}
	regencyDTO.Districts = districtsDTO

	cacheData := map[string]interface{}{
		"data": regencyDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching regency data to Redis: %v\n", err)
	}

	return regencyDTO, nil
}

func (s *wilayahIndonesiaService) GetAllDistricts(page, limit int) ([]dto.DistrictResponseDTO, error) {

	cacheKey := fmt.Sprintf("districts_page_%d_limit_%d", page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var districts []dto.DistrictResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				district, ok := item.(map[string]interface{})
				if ok {
					districts = append(districts, dto.DistrictResponseDTO{
						ID:        district["id"].(string),
						RegencyID: district["regency_id"].(string),
						Name:      district["name"].(string),
					})
				}
			}
			return districts, nil
		}
	}

	districts, total, err := s.WilayahRepo.FindAllDistricts(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch districts: %v", err)
	}

	var districtDTOs []dto.DistrictResponseDTO
	for _, district := range districts {
		districtDTOs = append(districtDTOs, dto.DistrictResponseDTO{
			ID:        district.ID,
			RegencyID: district.RegencyID,
			Name:      district.Name,
		})
	}

	cacheData := map[string]interface{}{
		"data":  districtDTOs,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching districts data to Redis: %v\n", err)
	}

	return districtDTOs, nil
}

func (s *wilayahIndonesiaService) GetDistrictByID(id string) (*dto.DistrictResponseDTO, error) {

	cacheKey := fmt.Sprintf("district:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var district dto.DistrictResponseDTO
		if data, ok := cachedData["data"].(map[string]interface{}); ok {
			district = dto.DistrictResponseDTO{
				ID:        data["id"].(string),
				RegencyID: data["regency_id"].(string),
				Name:      data["name"].(string),
			}
			return &district, nil
		}
	}

	district, err := s.WilayahRepo.FindDistrictByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %v", err)
	}

	districtDTO := &dto.DistrictResponseDTO{
		ID:        district.ID,
		RegencyID: district.RegencyID,
		Name:      district.Name,
	}

	villagesDTO := []dto.VillageResponseDTO{}
	for _, village := range district.Villages {
		villagesDTO = append(villagesDTO, dto.VillageResponseDTO{
			ID:         village.ID,
			DistrictID: village.DistrictID,
			Name:       village.Name,
		})
	}
	districtDTO.Villages = villagesDTO

	cacheData := map[string]interface{}{
		"data": districtDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching district data to Redis: %v\n", err)
	}

	return districtDTO, nil
}

func (s *wilayahIndonesiaService) GetAllVillages(page, limit int) ([]dto.VillageResponseDTO, error) {

	cacheKey := fmt.Sprintf("villages:%d:%d", page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var villages []dto.VillageResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				village, ok := item.(map[string]interface{})
				if ok {
					villages = append(villages, dto.VillageResponseDTO{
						ID:         village["id"].(string),
						DistrictID: village["district_id"].(string),
						Name:       village["name"].(string),
					})
				}
			}
			return villages, nil
		}
	}

	villages, total, err := s.WilayahRepo.FindAllVillages(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch villages: %v", err)
	}

	var villageDTOs []dto.VillageResponseDTO
	for _, village := range villages {
		villageDTOs = append(villageDTOs, dto.VillageResponseDTO{
			ID:         village.ID,
			DistrictID: village.DistrictID,
			Name:       village.Name,
		})
	}

	cacheData := map[string]interface{}{
		"data":  villageDTOs,
		"total": total,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching village data to Redis: %v\n", err)
	}

	return villageDTOs, nil
}

func (s *wilayahIndonesiaService) GetVillageByID(id string) (*dto.VillageResponseDTO, error) {

	cacheKey := fmt.Sprintf("village:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var villageDTO dto.VillageResponseDTO
		if data, ok := cachedData["data"].(map[string]interface{}); ok {
			villageDTO = dto.VillageResponseDTO{
				ID:         data["id"].(string),
				DistrictID: data["district_id"].(string),
				Name:       data["name"].(string),
			}
			return &villageDTO, nil
		}
	}

	village, err := s.WilayahRepo.FindVillageByID(id)
	if err != nil {
		return nil, fmt.Errorf("village not found: %v", err)
	}

	villageDTO := &dto.VillageResponseDTO{
		ID:         village.ID,
		DistrictID: village.DistrictID,
		Name:       village.Name,
	}

	cacheData := map[string]interface{}{
		"data": villageDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching village data to Redis: %v\n", err)
	}

	return villageDTO, nil
}
