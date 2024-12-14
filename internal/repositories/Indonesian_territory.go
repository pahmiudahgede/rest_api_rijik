package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProvinces() ([]domain.Province, error) {
	records, err := utils.ReadCSV("public/document/provinces.csv")
	if err != nil {
		return nil, err
	}

	var provinces []domain.Province
	for _, record := range records {
		province := domain.Province{
			ID:   record[0],
			Name: record[1],
		}
		provinces = append(provinces, province)
	}

	return provinces, nil
}

func GetRegencies() ([]domain.Regency, error) {
	records, err := utils.ReadCSV("public/document/regencies.csv")
	if err != nil {
		return nil, err
	}

	var regencies []domain.Regency
	for _, record := range records {
		regency := domain.Regency{
			ID:         record[0],
			ProvinceID: record[1],
			Name:       record[2],
		}
		regencies = append(regencies, regency)
	}

	return regencies, nil
}

func GetDistricts() ([]domain.District, error) {
	records, err := utils.ReadCSV("public/document/districts.csv")
	if err != nil {
		return nil, err
	}

	var districts []domain.District
	for _, record := range records {
		district := domain.District{
			ID:        record[0],
			RegencyID: record[1],
			Name:      record[2],
		}
		districts = append(districts, district)
	}

	return districts, nil
}

func GetVillages() ([]domain.Village, error) {
	records, err := utils.ReadCSV("public/document/villages.csv")
	if err != nil {
		return nil, err
	}

	var villages []domain.Village
	for _, record := range records {
		village := domain.Village{
			ID:         record[0],
			DistrictID: record[1],
			Name:       record[2],
		}
		villages = append(villages, village)
	}

	return villages, nil
}

func GetProvinceByID(id string) (domain.Province, error) {
	provinces, err := GetProvinces()
	if err != nil {
		return domain.Province{}, err
	}

	for _, province := range provinces {
		if province.ID == id {

			regencies, err := GetRegenciesByProvinceID(id)
			if err != nil {
				return domain.Province{}, err
			}

			province.ListRegency = regencies
			return province, nil
		}
	}
	return domain.Province{}, errors.New("province not found")
}

func GetRegencyByID(id string) (domain.Regency, error) {
	regencies, err := GetRegencies()
	if err != nil {
		return domain.Regency{}, err
	}

	for _, regency := range regencies {
		if regency.ID == id {

			districts, err := GetDistrictsByRegencyID(id)
			if err != nil {
				return domain.Regency{}, err
			}

			regency.ListDistrict = districts
			return regency, nil
		}
	}
	return domain.Regency{}, errors.New("regency not found")
}

func GetDistrictByID(id string) (domain.District, error) {
	districts, err := GetDistricts()
	if err != nil {
		return domain.District{}, err
	}

	for _, district := range districts {
		if district.ID == id {

			villages, err := GetVillagesByDistrictID(id)
			if err != nil {
				return domain.District{}, err
			}

			district.ListVillage = villages
			return district, nil
		}
	}
	return domain.District{}, errors.New("district not found")
}

func GetVillageByID(id string) (domain.Village, error) {
	villages, err := GetVillages()
	if err != nil {
		return domain.Village{}, err
	}

	for _, village := range villages {
		if village.ID == id {
			return village, nil
		}
	}
	return domain.Village{}, errors.New("village not found")
}

func GetRegenciesByProvinceID(provinceID string) ([]domain.Regency, error) {
	regencies, err := GetRegencies()
	if err != nil {
		return nil, err
	}

	var result []domain.Regency
	for _, regency := range regencies {
		if regency.ProvinceID == provinceID {
			result = append(result, regency)
		}
	}
	return result, nil
}

func GetDistrictsByRegencyID(regencyID string) ([]domain.District, error) {
	districts, err := GetDistricts()
	if err != nil {
		return nil, err
	}

	var result []domain.District
	for _, district := range districts {
		if district.RegencyID == regencyID {
			result = append(result, district)
		}
	}
	return result, nil
}

func GetVillagesByDistrictID(districtID string) ([]domain.Village, error) {
	villages, err := GetVillages()
	if err != nil {
		return nil, err
	}

	var result []domain.Village
	for _, village := range villages {
		if village.DistrictID == districtID {
			result = append(result, village)
		}
	}
	return result, nil
}
