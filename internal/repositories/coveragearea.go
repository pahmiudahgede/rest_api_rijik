package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetCoverageAreas() ([]domain.CoverageArea, error) {
	var coverageAreas []domain.CoverageArea
	if err := config.DB.Find(&coverageAreas).Error; err != nil {
		return nil, err
	}
	return coverageAreas, nil
}

func GetCoverageAreaByID(id string) (domain.CoverageArea, error) {
	var coverageArea domain.CoverageArea
	if err := config.DB.Where("id = ?", id).First(&coverageArea).Error; err != nil {
		return coverageArea, err
	}
	return coverageArea, nil
}

func GetCoverageAreaByDistrictID(id string) (domain.CoverageDistric, error) {
	var coverageDistric domain.CoverageDistric
	if err := config.DB.Where("id = ?", id).First(&coverageDistric).Error; err != nil {
		return coverageDistric, err
	}
	return coverageDistric, nil
}

func GetCoverageDistricsByCoverageAreaID(areaID string) ([]domain.CoverageDistric, error) {
	var districts []domain.CoverageDistric
	if err := config.DB.Where("coverage_area_id = ?", areaID).Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

func GetSubdistrictsByCoverageDistrictID(districtID string) ([]domain.CoverageSubdistrict, error) {
	var subdistricts []domain.CoverageSubdistrict
	if err := config.DB.Where("coverage_district_id = ?", districtID).Find(&subdistricts).Error; err != nil {
		return nil, err
	}
	return subdistricts, nil
}

func CreateCoverageArea(coverageArea *domain.CoverageArea) error {
	if err := config.DB.Create(&coverageArea).Error; err != nil {
		return err
	}
	return nil
}

func CreateCoverageDistrict(coverageDistrict *domain.CoverageDistric) error {
	if err := config.DB.Create(&coverageDistrict).Error; err != nil {
		return err
	}
	return nil
}

func CreateCoverageSubdistrict(coverageSubdistrict *domain.CoverageSubdistrict) error {
	if err := config.DB.Create(&coverageSubdistrict).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCoverageArea(id string, coverageArea domain.CoverageArea) (domain.CoverageArea, error) {
	var existingCoverageArea domain.CoverageArea
	if err := config.DB.Where("id = ?", id).First(&existingCoverageArea).Error; err != nil {
		return existingCoverageArea, err
	}

	existingCoverageArea.Province = coverageArea.Province
	if err := config.DB.Save(&existingCoverageArea).Error; err != nil {
		return existingCoverageArea, err
	}

	return existingCoverageArea, nil
}

func UpdateCoverageDistrict(id string, coverageDistrict domain.CoverageDistric) (domain.CoverageDistric, error) {
	var existingCoverageDistrict domain.CoverageDistric
	if err := config.DB.Where("id = ?", id).First(&existingCoverageDistrict).Error; err != nil {
		return existingCoverageDistrict, err
	}

	existingCoverageDistrict.District = coverageDistrict.District
	if err := config.DB.Save(&existingCoverageDistrict).Error; err != nil {
		return existingCoverageDistrict, err
	}

	return existingCoverageDistrict, nil
}

func UpdateCoverageSubdistrict(id string, coverageSubdistrict domain.CoverageSubdistrict) (domain.CoverageSubdistrict, error) {
	var existingCoverageSubdistrict domain.CoverageSubdistrict
	if err := config.DB.Where("id = ?", id).First(&existingCoverageSubdistrict).Error; err != nil {
		return existingCoverageSubdistrict, err
	}

	existingCoverageSubdistrict.Subdistrict = coverageSubdistrict.Subdistrict
	if err := config.DB.Save(&existingCoverageSubdistrict).Error; err != nil {
		return existingCoverageSubdistrict, err
	}

	return existingCoverageSubdistrict, nil
}

func DeleteCoverageArea(id string) error {
	var coverageArea domain.CoverageArea
	if err := config.DB.Where("id = ?", id).First(&coverageArea).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&coverageArea).Error; err != nil {
		return err
	}

	return nil
}

func DeleteCoverageDistrict(id string) error {
	var coverageDistrict domain.CoverageDistric
	if err := config.DB.Where("id = ?", id).First(&coverageDistrict).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&coverageDistrict).Error; err != nil {
		return err
	}

	return nil
}

func DeleteCoverageSubdistrict(id string) error {
	var coverageSubdistrict domain.CoverageSubdistrict
	if err := config.DB.Where("id = ?", id).First(&coverageSubdistrict).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&coverageSubdistrict).Error; err != nil {
		return err
	}

	return nil
}
