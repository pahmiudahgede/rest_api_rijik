package services

import (
	"fmt"
	"log"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type CoverageAreaService interface {
	CreateCoverageArea(request dto.RequestCoverageArea) (*dto.ResponseCoverageArea, error)
	GetCoverageAreaByID(id string) (*dto.ResponseCoverageArea, error)
	GetAllCoverageAreas() ([]dto.ResponseCoverageArea, error)
	UpdateCoverageArea(id string, request dto.RequestCoverageArea) (*dto.ResponseCoverageArea, error)
	DeleteCoverageArea(id string) error
}

type coverageAreaService struct {
	repo        repositories.CoverageAreaRepository
	WilayahRepo repositories.WilayahIndonesiaRepository
}

func NewCoverageAreaService(repo repositories.CoverageAreaRepository, WilayahRepo repositories.WilayahIndonesiaRepository) CoverageAreaService {
	return &coverageAreaService{repo: repo, WilayahRepo: WilayahRepo}
}

func ConvertCoverageAreaToResponse(coverage *model.CoverageArea) *dto.ResponseCoverageArea {
	createdAt, _ := utils.FormatDateToIndonesianFormat(coverage.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(coverage.UpdatedAt)

	return &dto.ResponseCoverageArea{
		ID:        coverage.ID,
		Province:  coverage.Province,
		Regency:   coverage.Regency,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (s *coverageAreaService) CreateCoverageArea(request dto.RequestCoverageArea) (*dto.ResponseCoverageArea, error) {
	errors, valid := request.ValidateCoverageArea()
	if !valid {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	province, _, err := s.WilayahRepo.FindProvinceByID(request.Province, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid province_id")
	}

	regency, _, err := s.WilayahRepo.FindRegencyByID(request.Regency, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid regency_id")
	}

	existingCoverage, err := s.repo.FindCoverageByProvinceAndRegency(province.Name, regency.Name)
	if err == nil && existingCoverage != nil {
		return nil, fmt.Errorf("coverage area with province %s and regency %s already exists", province.Name, regency.Name)
	}

	coverage := model.CoverageArea{
		Province: province.Name,
		Regency:  regency.Name,
	}

	if err := s.repo.CreateCoverage(&coverage); err != nil {
		return nil, fmt.Errorf("failed to create coverage area: %v", err)
	}

	response := ConvertCoverageAreaToResponse(&coverage)

	return response, nil
}

func (s *coverageAreaService) GetCoverageAreaByID(id string) (*dto.ResponseCoverageArea, error) {
	coverage, err := s.repo.FindCoverageById(id)
	if err != nil {
		return nil, err
	}

	response := ConvertCoverageAreaToResponse(coverage)

	return response, nil
}

func (s *coverageAreaService) GetAllCoverageAreas() ([]dto.ResponseCoverageArea, error) {
	coverageAreas, err := s.repo.FindAllCoverage()
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseCoverageArea
	for _, coverage := range coverageAreas {

		response = append(response, *ConvertCoverageAreaToResponse(&coverage))
	}

	return response, nil
}

func (s *coverageAreaService) UpdateCoverageArea(id string, request dto.RequestCoverageArea) (*dto.ResponseCoverageArea, error) {

	errors, valid := request.ValidateCoverageArea()
	if !valid {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	coverage, err := s.repo.FindCoverageById(id)
	if err != nil {
		return nil, fmt.Errorf("coverage area with ID %s not found: %v", id, err)
	}

	province, _, err := s.WilayahRepo.FindProvinceByID(request.Province, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid province_id")
	}

	regency, _, err := s.WilayahRepo.FindRegencyByID(request.Regency, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid regency_id")
	}

	existingCoverage, err := s.repo.FindCoverageByProvinceAndRegency(province.Name, regency.Name)
	if err == nil && existingCoverage != nil {
		return nil, fmt.Errorf("coverage area with province %s and regency %s already exists", province.Name, regency.Name)
	}

	coverage.Province = province.Name
	coverage.Regency = regency.Name

	if err := s.repo.UpdateCoverage(id, coverage); err != nil {
		return nil, fmt.Errorf("failed to update coverage area: %v", err)
	}

	response := ConvertCoverageAreaToResponse(coverage)

	return response, nil
}

func (s *coverageAreaService) DeleteCoverageArea(id string) error {

	coverage, err := s.repo.FindCoverageById(id)
	if err != nil {
		return fmt.Errorf("coverage area with ID %s not found: %v", id, err)
	}

	if err := s.repo.DeleteCoverage(id); err != nil {
		return fmt.Errorf("failed to delete coverage area: %v", err)
	}

	log.Printf("Coverage area with ID %s successfully deleted", coverage.ID)
	return nil
}
