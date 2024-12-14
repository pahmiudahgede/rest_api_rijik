package repositories

import (
	"encoding/csv"
	"os"

	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetProvinces() ([]domain.Province, error) {
	file, err := os.Open("public/document/provinces.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
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
	file, err := os.Open("public/document/regencies.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
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
	file, err := os.Open("public/document/districts.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
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
	file, err := os.Open("public/document/villages.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
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
