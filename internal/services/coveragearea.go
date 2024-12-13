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
