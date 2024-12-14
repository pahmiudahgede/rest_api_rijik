package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetProvinces() ([]domain.Province, error) {
	provinces, err := repositories.GetProvinces()
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func GetRegencies() ([]domain.Regency, error) {
	regencies, err := repositories.GetRegencies()
	if err != nil {
		return nil, err
	}
	return regencies, nil
}

func GetDistricts() ([]domain.District, error) {
	districts, err := repositories.GetDistricts()
	if err != nil {
		return nil, err
	}
	return districts, nil
}

func GetVillages() ([]domain.Village, error) {
	villages, err := repositories.GetVillages()
	if err != nil {
		return nil, err
	}
	return villages, nil
}

func GetProvinceByID(id string) (domain.Province, error) {
	province, err := repositories.GetProvinceByID(id)
	if err != nil {
		return domain.Province{}, err
	}
	return province, nil
}

func GetRegencyByID(id string) (domain.Regency, error) {
	regency, err := repositories.GetRegencyByID(id)
	if err != nil {
		return domain.Regency{}, err
	}
	return regency, nil
}

func GetDistrictByID(id string) (domain.District, error) {
	district, err := repositories.GetDistrictByID(id)
	if err != nil {
		return domain.District{}, err
	}
	return district, nil
}

func GetVillageByID(id string) (domain.Village, error) {
	village, err := repositories.GetVillageByID(id)
	if err != nil {
		return domain.Village{}, err
	}
	return village, nil
}
