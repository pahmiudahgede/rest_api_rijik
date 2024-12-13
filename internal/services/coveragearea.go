package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetCoverageAreas() ([]domain.CoverageArea, error) {
	return repositories.GetCoverageAreas()
}

func GetCoverageDistricsByCoverageAreaID(areaID string) ([]domain.CoverageDistric, error) {
	return repositories.GetCoverageDistricsByCoverageAreaID(areaID)
}

func GetCoverageAreaByDistrictID(id string) (domain.CoverageDistric, error) {
	return repositories.GetCoverageAreaByDistrictID(id)
}

func GetCoverageAreaByID(id string) (domain.CoverageArea, error) {
	return repositories.GetCoverageAreaByID(id)
}

func GetSubdistrictsByCoverageDistrictID(districtID string) ([]domain.CoverageSubdistrict, error) {
	return repositories.GetSubdistrictsByCoverageDistrictID(districtID)
}
