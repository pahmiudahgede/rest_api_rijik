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