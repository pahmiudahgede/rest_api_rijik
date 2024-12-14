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
