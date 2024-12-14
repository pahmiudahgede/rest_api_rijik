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

// GetRegencies retrieves a list of regencies
func GetRegencies() ([]domain.Regency, error) {
	regencies, err := repositories.GetRegencies()
	if err != nil {
		return nil, err
	}
	return regencies, nil
}

// GetDistricts retrieves a list of districts
func GetDistricts() ([]domain.District, error) {
	districts, err := repositories.GetDistricts()
	if err != nil {
		return nil, err
	}
	return districts, nil
}

// GetVillages retrieves a list of villages
func GetVillages() ([]domain.Village, error) {
	villages, err := repositories.GetVillages()
	if err != nil {
		return nil, err
	}
	return villages, nil
}

func GetProvinceByID(id string) (domain.Province, []domain.Regency, error) {
	province, err := repositories.FindProvinceByID(id)
	if err != nil {
		return domain.Province{}, nil, err
	}

	regencies, err := repositories.GetRegencies()
	if err != nil {
		return domain.Province{}, nil, err
	}

	var listRegency []domain.Regency
	for _, regency := range regencies {
		if regency.ProvinceID == province.ID {
			listRegency = append(listRegency, regency)
		}
	}

	return province, listRegency, nil
}

func GetRegencyByID(id string) (domain.Regency, []domain.District, error) {
	regency, err := repositories.FindRegencyByID(id)
	if err != nil {
		return domain.Regency{}, nil, err
	}

	districts, err := repositories.GetDistricts()
	if err != nil {
		return domain.Regency{}, nil, err
	}

	var listDistrict []domain.District
	for _, district := range districts {
		if district.RegencyID == regency.ID {
			listDistrict = append(listDistrict, district)
		}
	}

	return regency, listDistrict, nil
}

func GetDistrictByID(id string) (domain.District, []domain.Village, error) {
	district, err := repositories.FindDistrictByID(id)
	if err != nil {
		return domain.District{}, nil, err
	}

	villages, err := repositories.GetVillages()
	if err != nil {
		return domain.District{}, nil, err
	}

	var listVillage []domain.Village
	for _, village := range villages {
		if village.DistrictID == district.ID {
			listVillage = append(listVillage, village)
		}
	}

	return district, listVillage, nil
}

func GetVillageByID(id string) (domain.Village, error) {
	village, err := repositories.FindVillageByID(id)
	if err != nil {
		return domain.Village{}, err
	}
	return village, nil
}
