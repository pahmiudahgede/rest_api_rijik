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

// GetRegencies reads the regencies data from CSV and returns a slice of Regency
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

// GetDistricts reads the districts data from CSV and returns a slice of District
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

// GetVillages reads the villages data from CSV and returns a slice of Village
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


func FindProvinceByID(id string) (domain.Province, error) {
	provinces, err := GetProvinces()
	if err != nil {
		return domain.Province{}, err
	}

	for _, province := range provinces {
		if province.ID == id {
			return province, nil
		}
	}
	return domain.Province{}, errors.New("province not found")
}

func FindRegencyByID(id string) (domain.Regency, error) {
	regencies, err := GetRegencies()
	if err != nil {
		return domain.Regency{}, err
	}

	for _, regency := range regencies {
		if regency.ID == id {
			return regency, nil
		}
	}
	return domain.Regency{}, errors.New("regency not found")
}

func FindDistrictByID(id string) (domain.District, error) {
	districts, err := GetDistricts()
	if err != nil {
		return domain.District{}, err
	}

	for _, district := range districts {
		if district.ID == id {
			return district, nil
		}
	}
	return domain.District{}, errors.New("district not found")
}

func FindVillageByID(id string) (domain.Village, error) {
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
