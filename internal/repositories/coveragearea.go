package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetCoverageAreas() ([]domain.CoverageArea, error) {
	var coverageAreas []domain.CoverageArea
	if err := config.DB.Preload("Details").Find(&coverageAreas).Error; err != nil {
		return nil, err
	}
	return coverageAreas, nil
}

func GetCoverageAreaByID(id string) (domain.CoverageArea, error) {
	var coverageArea domain.CoverageArea
	if err := config.DB.Preload("Details.LocationSpecific").Where("id = ?", id).First(&coverageArea).Error; err != nil {
		return coverageArea, err
	}
	return coverageArea, nil
}

func GetCoverageAreaByDistrictID(id string) (domain.CoverageDetail, error) {
	var coverageDetail domain.CoverageDetail
	if err := config.DB.Preload("LocationSpecific").Where("id = ?", id).First(&coverageDetail).Error; err != nil {
		return coverageDetail, err
	}
	return coverageDetail, nil
}
