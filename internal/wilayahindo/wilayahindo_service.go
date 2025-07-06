package wilayahindo

import (
	"context"
	"fmt"
	"time"

	"rijig/model"
	"rijig/utils"
)

type WilayahIndonesiaService interface {
	ImportDataFromCSV(ctx context.Context) error

	GetAllProvinces(ctx context.Context, page, limit int) ([]ProvinceResponseDTO, int, error)
	GetProvinceByID(ctx context.Context, id string, page, limit int) (*ProvinceResponseDTO, int, error)

	GetAllRegencies(ctx context.Context, page, limit int) ([]RegencyResponseDTO, int, error)
	GetRegencyByID(ctx context.Context, id string, page, limit int) (*RegencyResponseDTO, int, error)

	GetAllDistricts(ctx context.Context, page, limit int) ([]DistrictResponseDTO, int, error)
	GetDistrictByID(ctx context.Context, id string, page, limit int) (*DistrictResponseDTO, int, error)

	GetAllVillages(ctx context.Context, page, limit int) ([]VillageResponseDTO, int, error)
	GetVillageByID(ctx context.Context, id string) (*VillageResponseDTO, error)
}

type wilayahIndonesiaService struct {
	WilayahRepo WilayahIndonesiaRepository
}

func NewWilayahIndonesiaService(wilayahRepo WilayahIndonesiaRepository) WilayahIndonesiaService {
	return &wilayahIndonesiaService{WilayahRepo: wilayahRepo}
}

func (s *wilayahIndonesiaService) ImportDataFromCSV(ctx context.Context) error {

	provinces, err := utils.ReadCSV("public/document/provinces.csv")
	if err != nil {
		return fmt.Errorf("failed to read provinces CSV: %w", err)
	}

	var provinceList []model.Province
	for _, record := range provinces[1:] {
		if len(record) >= 2 {
			province := model.Province{
				ID:   record[0],
				Name: record[1],
			}
			provinceList = append(provinceList, province)
		}
	}

	if err := s.WilayahRepo.ImportProvinces(ctx, provinceList); err != nil {
		return fmt.Errorf("failed to import provinces: %w", err)
	}

	regencies, err := utils.ReadCSV("public/document/regencies.csv")
	if err != nil {
		return fmt.Errorf("failed to read regencies CSV: %w", err)
	}

	var regencyList []model.Regency
	for _, record := range regencies[1:] {
		if len(record) >= 3 {
			regency := model.Regency{
				ID:         record[0],
				ProvinceID: record[1],
				Name:       record[2],
			}
			regencyList = append(regencyList, regency)
		}
	}

	if err := s.WilayahRepo.ImportRegencies(ctx, regencyList); err != nil {
		return fmt.Errorf("failed to import regencies: %w", err)
	}

	districts, err := utils.ReadCSV("public/document/districts.csv")
	if err != nil {
		return fmt.Errorf("failed to read districts CSV: %w", err)
	}

	var districtList []model.District
	for _, record := range districts[1:] {
		if len(record) >= 3 {
			district := model.District{
				ID:        record[0],
				RegencyID: record[1],
				Name:      record[2],
			}
			districtList = append(districtList, district)
		}
	}

	if err := s.WilayahRepo.ImportDistricts(ctx, districtList); err != nil {
		return fmt.Errorf("failed to import districts: %w", err)
	}

	villages, err := utils.ReadCSV("public/document/villages.csv")
	if err != nil {
		return fmt.Errorf("failed to read villages CSV: %w", err)
	}

	var villageList []model.Village
	for _, record := range villages[1:] {
		if len(record) >= 3 {
			village := model.Village{
				ID:         record[0],
				DistrictID: record[1],
				Name:       record[2],
			}
			villageList = append(villageList, village)
		}
	}

	if err := s.WilayahRepo.ImportVillages(ctx, villageList); err != nil {
		return fmt.Errorf("failed to import villages: %w", err)
	}

	return nil
}

func (s *wilayahIndonesiaService) GetAllProvinces(ctx context.Context, page, limit int) ([]ProvinceResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("provinces_page:%d_limit:%d", page, limit)

	var cachedResponse struct {
		Data  []ProvinceResponseDTO `json:"data"`
		Total int                       `json:"total"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return cachedResponse.Data, cachedResponse.Total, nil
	}

	provinces, total, err := s.WilayahRepo.FindAllProvinces(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch provinces: %w", err)
	}

	provinceDTOs := make([]ProvinceResponseDTO, len(provinces))
	for i, province := range provinces {
		provinceDTOs[i] = ProvinceResponseDTO{
			ID:   province.ID,
			Name: province.Name,
		}
	}

	cacheData := struct {
		Data  []ProvinceResponseDTO `json:"data"`
		Total int                       `json:"total"`
	}{
		Data:  provinceDTOs,
		Total: total,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching provinces data: %v\n", err)
	}

	return provinceDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetProvinceByID(ctx context.Context, id string, page, limit int) (*ProvinceResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("province:%s_page:%d_limit:%d", id, page, limit)

	var cachedResponse struct {
		Data           ProvinceResponseDTO `json:"data"`
		TotalRegencies int                     `json:"total_regencies"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse.Data, cachedResponse.TotalRegencies, nil
	}

	province, totalRegencies, err := s.WilayahRepo.FindProvinceByID(ctx, id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	provinceDTO := ProvinceResponseDTO{
		ID:   province.ID,
		Name: province.Name,
	}

	regencyDTOs := make([]RegencyResponseDTO, len(province.Regencies))
	for i, regency := range province.Regencies {
		regencyDTOs[i] = RegencyResponseDTO{
			ID:         regency.ID,
			ProvinceID: regency.ProvinceID,
			Name:       regency.Name,
		}
	}
	provinceDTO.Regencies = regencyDTOs

	cacheData := struct {
		Data           ProvinceResponseDTO `json:"data"`
		TotalRegencies int                     `json:"total_regencies"`
	}{
		Data:           provinceDTO,
		TotalRegencies: totalRegencies,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching province data: %v\n", err)
	}

	return &provinceDTO, totalRegencies, nil
}

func (s *wilayahIndonesiaService) GetAllRegencies(ctx context.Context, page, limit int) ([]RegencyResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("regencies_page:%d_limit:%d", page, limit)

	var cachedResponse struct {
		Data  []RegencyResponseDTO `json:"data"`
		Total int                      `json:"total"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return cachedResponse.Data, cachedResponse.Total, nil
	}

	regencies, total, err := s.WilayahRepo.FindAllRegencies(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch regencies: %w", err)
	}

	regencyDTOs := make([]RegencyResponseDTO, len(regencies))
	for i, regency := range regencies {
		regencyDTOs[i] = RegencyResponseDTO{
			ID:         regency.ID,
			ProvinceID: regency.ProvinceID,
			Name:       regency.Name,
		}
	}

	cacheData := struct {
		Data  []RegencyResponseDTO `json:"data"`
		Total int                      `json:"total"`
	}{
		Data:  regencyDTOs,
		Total: total,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching regencies data: %v\n", err)
	}

	return regencyDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetRegencyByID(ctx context.Context, id string, page, limit int) (*RegencyResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("regency:%s_page:%d_limit:%d", id, page, limit)

	var cachedResponse struct {
		Data           RegencyResponseDTO `json:"data"`
		TotalDistricts int                    `json:"total_districts"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse.Data, cachedResponse.TotalDistricts, nil
	}

	regency, totalDistricts, err := s.WilayahRepo.FindRegencyByID(ctx, id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	regencyDTO := RegencyResponseDTO{
		ID:         regency.ID,
		ProvinceID: regency.ProvinceID,
		Name:       regency.Name,
	}

	districtDTOs := make([]DistrictResponseDTO, len(regency.Districts))
	for i, district := range regency.Districts {
		districtDTOs[i] = DistrictResponseDTO{
			ID:        district.ID,
			RegencyID: district.RegencyID,
			Name:      district.Name,
		}
	}
	regencyDTO.Districts = districtDTOs

	cacheData := struct {
		Data           RegencyResponseDTO `json:"data"`
		TotalDistricts int                    `json:"total_districts"`
	}{
		Data:           regencyDTO,
		TotalDistricts: totalDistricts,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching regency data: %v\n", err)
	}

	return &regencyDTO, totalDistricts, nil
}

func (s *wilayahIndonesiaService) GetAllDistricts(ctx context.Context, page, limit int) ([]DistrictResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("districts_page:%d_limit:%d", page, limit)

	var cachedResponse struct {
		Data  []DistrictResponseDTO `json:"data"`
		Total int                       `json:"total"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return cachedResponse.Data, cachedResponse.Total, nil
	}

	districts, total, err := s.WilayahRepo.FindAllDistricts(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch districts: %w", err)
	}

	districtDTOs := make([]DistrictResponseDTO, len(districts))
	for i, district := range districts {
		districtDTOs[i] = DistrictResponseDTO{
			ID:        district.ID,
			RegencyID: district.RegencyID,
			Name:      district.Name,
		}
	}

	cacheData := struct {
		Data  []DistrictResponseDTO `json:"data"`
		Total int                       `json:"total"`
	}{
		Data:  districtDTOs,
		Total: total,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching districts data: %v\n", err)
	}

	return districtDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetDistrictByID(ctx context.Context, id string, page, limit int) (*DistrictResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("district:%s_page:%d_limit:%d", id, page, limit)

	var cachedResponse struct {
		Data          DistrictResponseDTO `json:"data"`
		TotalVillages int                     `json:"total_villages"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse.Data, cachedResponse.TotalVillages, nil
	}

	district, totalVillages, err := s.WilayahRepo.FindDistrictByID(ctx, id, page, limit)
	if err != nil {
		return nil, 0, err
	}

	districtDTO := DistrictResponseDTO{
		ID:        district.ID,
		RegencyID: district.RegencyID,
		Name:      district.Name,
	}

	villageDTOs := make([]VillageResponseDTO, len(district.Villages))
	for i, village := range district.Villages {
		villageDTOs[i] = VillageResponseDTO{
			ID:         village.ID,
			DistrictID: village.DistrictID,
			Name:       village.Name,
		}
	}
	districtDTO.Villages = villageDTOs

	cacheData := struct {
		Data          DistrictResponseDTO `json:"data"`
		TotalVillages int                     `json:"total_villages"`
	}{
		Data:          districtDTO,
		TotalVillages: totalVillages,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching district data: %v\n", err)
	}

	return &districtDTO, totalVillages, nil
}

func (s *wilayahIndonesiaService) GetAllVillages(ctx context.Context, page, limit int) ([]VillageResponseDTO, int, error) {
	cacheKey := fmt.Sprintf("villages_page:%d_limit:%d", page, limit)

	var cachedResponse struct {
		Data  []VillageResponseDTO `json:"data"`
		Total int                      `json:"total"`
	}

	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return cachedResponse.Data, cachedResponse.Total, nil
	}

	villages, total, err := s.WilayahRepo.FindAllVillages(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch villages: %w", err)
	}

	villageDTOs := make([]VillageResponseDTO, len(villages))
	for i, village := range villages {
		villageDTOs[i] = VillageResponseDTO{
			ID:         village.ID,
			DistrictID: village.DistrictID,
			Name:       village.Name,
		}
	}

	cacheData := struct {
		Data  []VillageResponseDTO `json:"data"`
		Total int                      `json:"total"`
	}{
		Data:  villageDTOs,
		Total: total,
	}

	if err := utils.SetCache(cacheKey, cacheData, 24*time.Hour); err != nil {
		fmt.Printf("Error caching villages data: %v\n", err)
	}

	return villageDTOs, total, nil
}

func (s *wilayahIndonesiaService) GetVillageByID(ctx context.Context, id string) (*VillageResponseDTO, error) {
	cacheKey := fmt.Sprintf("village:%s", id)

	var cachedResponse VillageResponseDTO
	if err := utils.GetCache(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	village, err := s.WilayahRepo.FindVillageByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("village not found: %w", err)
	}

	villageResponse := &VillageResponseDTO{
		ID:         village.ID,
		DistrictID: village.DistrictID,
		Name:       village.Name,
	}

	if err := utils.SetCache(cacheKey, villageResponse, 24*time.Hour); err != nil {
		fmt.Printf("Error caching village data: %v\n", err)
	}

	return villageResponse, nil
}
