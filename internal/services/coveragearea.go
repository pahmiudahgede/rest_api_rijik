package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetCoverageAreas() ([]domain.CoverageArea, error) {
	return repositories.GetCoverageAreas()
}

func GetCoverageAreaByID(id string) (domain.CoverageArea, error) {
	return repositories.GetCoverageAreaByID(id)
}

func GetCoverageAreaByDistrictID(id string) (domain.CoverageDistric, error) {
	return repositories.GetCoverageAreaByDistrictID(id)
}

func GetCoverageDistricsByCoverageAreaID(areaID string) ([]domain.CoverageDistric, error) {
	return repositories.GetCoverageDistricsByCoverageAreaID(areaID)
}

func GetSubdistrictsByCoverageDistrictID(districtID string) ([]domain.CoverageSubdistrict, error) {
	return repositories.GetSubdistrictsByCoverageDistrictID(districtID)
}

func CreateCoverageArea(province string) (*domain.CoverageArea, error) {
	coverageArea := &domain.CoverageArea{
		Province: province,
	}

	if err := repositories.CreateCoverageArea(coverageArea); err != nil {
		return nil, err
	}

	return coverageArea, nil
}

func CreateCoverageDistrict(coverageAreaID, district string) (*domain.CoverageDistric, error) {
	coverageDistrict := &domain.CoverageDistric{
		CoverageAreaID: coverageAreaID,
		District:       district,
	}

	if err := repositories.CreateCoverageDistrict(coverageDistrict); err != nil {
		return nil, err
	}

	return coverageDistrict, nil
}

func CreateCoverageSubdistrict(coverageAreaID, coverageDistrictId, subdistrict string) (*domain.CoverageSubdistrict, error) {
	coverageSubdistrict := &domain.CoverageSubdistrict{
		CoverageAreaID:     coverageAreaID,
		CoverageDistrictId: coverageDistrictId,
		Subdistrict:        subdistrict,
	}

	if err := repositories.CreateCoverageSubdistrict(coverageSubdistrict); err != nil {
		return nil, err
	}

	return coverageSubdistrict, nil
}