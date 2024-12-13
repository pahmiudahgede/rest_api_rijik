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

func GetCoverageAreaByDistrictID(id string) (domain.CoverageDetail, error) {
	return repositories.GetCoverageAreaByDistrictID(id)
}

func CreateCoverageArea(province string) (domain.CoverageArea, error) {
	return repositories.CreateCoverageArea(province)
}

func CreateCoverageDetail(coverageAreaID, province, district string) (domain.CoverageDetail, error) {
	return repositories.CreateCoverageDetail(coverageAreaID, province, district)
}

func CreateLocationSpecific(coverageDetailID, subdistrict string) (domain.LocationSpecific, error) {
	return repositories.CreateLocationSpecific(coverageDetailID, subdistrict)
}