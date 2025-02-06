package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type WilayahIndonesiaService interface {
	ImportDataFromCSV() error

	GetAllProvinces(page, limit int) ([]dto.ProvinceResponseDTO, int, error)
	GetProvinceByID(id string, page, limit int) (*dto.ProvinceResponseDTO, int, error)

	GetAllRegencies(page, limit int) ([]dto.RegencyResponseDTO, int, error)
	GetRegencyByID(id string, page, limit int) (*dto.RegencyResponseDTO, int, error)

	GetAllDistricts(page, limit int) ([]dto.DistrictResponseDTO, int, error)
	GetDistrictByID(id string, page, limit int) (*dto.DistrictResponseDTO, int, error)

	GetAllVillages(page, limit int) ([]dto.VillageResponseDTO, int, error)
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

func (s *wilayahIndonesiaService) GetAllProvinces(page, limit int) ([]dto.ProvinceResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("provinces_page:%d_limit:%d", page, limit)
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
			total := int(cachedData["total"].(float64))
			return provinces, total, nil
		}
	}

	provinces, total, err := s.WilayahRepo.FindAllProvinces(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch provinces: %v", err)
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
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching provinces data: %v\n", err)
	}

	return provinceDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetProvinceByID(id string, page, limit int) (*dto.ProvinceResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("province:%s_page:%d_limit:%d", id, page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		var provinceDTO dto.ProvinceResponseDTO
		if data, ok := cachedData["data"].(string); ok {
			if err := json.Unmarshal([]byte(data), &provinceDTO); err == nil {

				totalRegencies, _ := strconv.Atoi(cachedData["total_regencies"].(string))
				return &provinceDTO, totalRegencies, nil
			}
		}
	}

	province, totalRegencies, err := s.WilayahRepo.FindProvinceByID(id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	provinceDTO := dto.ProvinceResponseDTO{
		ID:   province.ID,
		Name: province.Name,
	}

	var regencyDTOs []dto.RegencyResponseDTO
	for _, regency := range province.Regencies {
		regencyDTO := dto.RegencyResponseDTO{
			ID:         regency.ID,
			ProvinceID: regency.ProvinceID,
			Name:       regency.Name,
		}
		regencyDTOs = append(regencyDTOs, regencyDTO)
	}

	provinceDTO.Regencies = regencyDTOs

	cacheData := map[string]interface{}{
		"data":            provinceDTO,
		"total_regencies": strconv.Itoa(totalRegencies),
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching province data: %v\n", err)
	}

	return &provinceDTO, totalRegencies, nil
}

func (s *wilayahIndonesiaService) GetAllRegencies(page, limit int) ([]dto.RegencyResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("regencies_page:%d_limit:%d", page, limit)
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
			total := int(cachedData["total"].(float64))
			return regencies, total, nil
		}
	}

	regencies, total, err := s.WilayahRepo.FindAllRegencies(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch provinces: %v", err)
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
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching regencies data: %v\n", err)
	}

	return regencyDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetRegencyByID(id string, page, limit int) (*dto.RegencyResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("regency:%s_page:%d_limit:%d", id, page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		var regencyDTO dto.RegencyResponseDTO
		if data, ok := cachedData["data"].(string); ok {
			if err := json.Unmarshal([]byte(data), &regencyDTO); err == nil {

				totalDistrict, _ := strconv.Atoi(cachedData["total_regencies"].(string))
				return &regencyDTO, totalDistrict, nil
			}
		}
	}

	regency, totalDistrict, err := s.WilayahRepo.FindRegencyByID(id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	regencyDTO := dto.RegencyResponseDTO{
		ID:         regency.ID,
		ProvinceID: regency.ProvinceID,
		Name:       regency.Name,
	}

	var districtDTOs []dto.DistrictResponseDTO
	for _, regency := range regency.Districts {
		districtDTO := dto.DistrictResponseDTO{
			ID:        regency.ID,
			RegencyID: regency.RegencyID,
			Name:      regency.Name,
		}
		districtDTOs = append(districtDTOs, districtDTO)
	}

	regencyDTO.Districts = districtDTOs

	cacheData := map[string]interface{}{
		"data":            regencyDTO,
		"total_regencies": strconv.Itoa(totalDistrict),
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching province data: %v\n", err)
	}

	return &regencyDTO, totalDistrict, nil
}

func (s *wilayahIndonesiaService) GetAllDistricts(page, limit int) ([]dto.DistrictResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("district_page:%d_limit:%d", page, limit)
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
			total := int(cachedData["total"].(float64))
			return districts, total, nil
		}
	}

	districts, total, err := s.WilayahRepo.FindAllDistricts(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch districts: %v", err)
	}

	var districtsDTOs []dto.DistrictResponseDTO
	for _, district := range districts {
		districtsDTOs = append(districtsDTOs, dto.DistrictResponseDTO{
			ID:        district.ID,
			RegencyID: district.RegencyID,
			Name:      district.Name,
		})
	}

	cacheData := map[string]interface{}{
		"data":  districtsDTOs,
		"total": total,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching districts data: %v\n", err)
	}

	return districtsDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetDistrictByID(id string, page, limit int) (*dto.DistrictResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("district:%s_page:%d_limit:%d", id, page, limit)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		var districtDTO dto.DistrictResponseDTO
		if data, ok := cachedData["data"].(string); ok {
			if err := json.Unmarshal([]byte(data), &districtDTO); err == nil {

				totalVillage, _ := strconv.Atoi(cachedData["total_village"].(string))
				return &districtDTO, totalVillage, nil
			}
		}
	}

	district, totalVillages, err := s.WilayahRepo.FindDistrictByID(id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	districtDTO := dto.DistrictResponseDTO{
		ID:        district.ID,
		RegencyID: district.RegencyID,
		Name:      district.Name,
	}

	var villageDTOs []dto.VillageResponseDTO
	for _, village := range district.Villages {
		regencyDTO := dto.VillageResponseDTO{
			ID:         village.ID,
			DistrictID: village.DistrictID,
			Name:       village.Name,
		}
		villageDTOs = append(villageDTOs, regencyDTO)
	}

	districtDTO.Villages = villageDTOs

	cacheData := map[string]interface{}{
		"data":           districtDTO,
		"total_villages": strconv.Itoa(totalVillages),
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching province data: %v\n", err)
	}

	return &districtDTO, totalVillages, nil
}

func (s *wilayahIndonesiaService) GetAllVillages(page, limit int) ([]dto.VillageResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("villages_page:%d_limit:%d", page, limit)
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
			total := int(cachedData["total"].(float64))
			return villages, total, nil
		}
	}

	villages, total, err := s.WilayahRepo.FindAllVillages(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch villages: %v", err)
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
		fmt.Printf("Error caching villages data: %v\n", err)
	}

	return villageDTOs, total, nil
}
