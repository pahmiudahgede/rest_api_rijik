package repositories

import (
	"errors"

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
	if err := config.DB.Preload("Details").Where("id = ?", id).First(&coverageArea).Error; err != nil {
		return coverageArea, errors.New("coverage area not found")
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

func CreateCoverageArea(province string) (domain.CoverageArea, error) {
	coverageArea := domain.CoverageArea{
		Province: province,
	}

	if err := config.DB.Create(&coverageArea).Error; err != nil {
		return domain.CoverageArea{}, err
	}

	return coverageArea, nil
}

func CreateCoverageDetail(coverageAreaID, province, district string) (domain.CoverageDetail, error) {
	coverageDetail := domain.CoverageDetail{
		CoverageAreaID: coverageAreaID,
		District: district,
	}

	if err := config.DB.Create(&coverageDetail).Error; err != nil {
		return domain.CoverageDetail{}, err
	}

	return coverageDetail, nil
}

func CreateLocationSpecific(coverageDetailID, subdistrict string) (domain.LocationSpecific, error) {
	locationSpecific := domain.LocationSpecific{
		CoverageDetailID: coverageDetailID,
		Subdistrict:      subdistrict,
	}

	if err := config.DB.Create(&locationSpecific).Error; err != nil {
		return domain.LocationSpecific{}, err
	}

	return locationSpecific, nil
}
