package services

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type CompanyProfileService interface {
	CreateCompanyProfile(userID string, request *dto.RequestCompanyProfileDTO) (*dto.ResponseCompanyProfileDTO, error)
	GetCompanyProfileByID(id string) (*dto.ResponseCompanyProfileDTO, error)
	GetCompanyProfilesByUserID(userID string) ([]dto.ResponseCompanyProfileDTO, error)
	UpdateCompanyProfile(id string, request *dto.RequestCompanyProfileDTO) (*dto.ResponseCompanyProfileDTO, error)
	DeleteCompanyProfile(id string) error
}

type companyProfileService struct {
	companyProfileRepo repositories.CompanyProfileRepository
}

func NewCompanyProfileService(companyProfileRepo repositories.CompanyProfileRepository) CompanyProfileService {
	return &companyProfileService{
		companyProfileRepo: companyProfileRepo,
	}
}

func FormatResponseCompanyProfile(companyProfile *model.CompanyProfile) (*dto.ResponseCompanyProfileDTO, error) {

	createdAt, _ := utils.FormatDateToIndonesianFormat(companyProfile.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(companyProfile.UpdatedAt)

	responseDTO := &dto.ResponseCompanyProfileDTO{
		ID:                 companyProfile.ID,
		UserID:             companyProfile.UserID,
		CompanyName:        companyProfile.CompanyName,
		CompanyAddress:     companyProfile.CompanyAddress,
		CompanyPhone:       companyProfile.CompanyPhone,
		CompanyEmail:       companyProfile.CompanyEmail,
		CompanyLogo:        companyProfile.CompanyLogo,
		CompanyWebsite:     companyProfile.CompanyWebsite,
		TaxID:              companyProfile.TaxID,
		FoundedDate:        companyProfile.FoundedDate,
		CompanyType:        companyProfile.CompanyType,
		CompanyDescription: companyProfile.CompanyDescription,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	return responseDTO, nil
}

func (s *companyProfileService) CreateCompanyProfile(userID string, request *dto.RequestCompanyProfileDTO) (*dto.ResponseCompanyProfileDTO, error) {

	errors, valid := request.ValidateCompanyProfileInput()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	companyProfile := &model.CompanyProfile{
		UserID:             userID,
		CompanyName:        request.CompanyName,
		CompanyAddress:     request.CompanyAddress,
		CompanyPhone:       request.CompanyPhone,
		CompanyEmail:       request.CompanyEmail,
		CompanyLogo:        request.CompanyLogo,
		CompanyWebsite:     request.CompanyWebsite,
		TaxID:              request.TaxID,
		FoundedDate:        request.FoundedDate,
		CompanyType:        request.CompanyType,
		CompanyDescription: request.CompanyDescription,
	}

	createdCompanyProfile, err := s.companyProfileRepo.CreateCompanyProfile(companyProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create company profile: %v", err)
	}

	responseDTO, err := FormatResponseCompanyProfile(createdCompanyProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to format company profile response: %v", err)
	}

	return responseDTO, nil
}

func (s *companyProfileService) GetCompanyProfileByID(id string) (*dto.ResponseCompanyProfileDTO, error) {

	companyProfile, err := s.companyProfileRepo.GetCompanyProfileByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving company profile by ID: %v", err)
	}

	responseDTO, err := FormatResponseCompanyProfile(companyProfile)
	if err != nil {
		return nil, fmt.Errorf("error formatting company profile response: %v", err)
	}

	return responseDTO, nil
}

func (s *companyProfileService) GetCompanyProfilesByUserID(userID string) ([]dto.ResponseCompanyProfileDTO, error) {

	companyProfiles, err := s.companyProfileRepo.GetCompanyProfilesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving company profiles by userID: %v", err)
	}

	var responseDTOs []dto.ResponseCompanyProfileDTO
	for _, companyProfile := range companyProfiles {
		responseDTO, err := FormatResponseCompanyProfile(&companyProfile)
		if err != nil {
			return nil, fmt.Errorf("error formatting company profile response: %v", err)
		}
		responseDTOs = append(responseDTOs, *responseDTO)
	}

	return responseDTOs, nil
}

func (s *companyProfileService) UpdateCompanyProfile(id string, request *dto.RequestCompanyProfileDTO) (*dto.ResponseCompanyProfileDTO, error) {

	errors, valid := request.ValidateCompanyProfileInput()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	companyProfile := &model.CompanyProfile{
		CompanyName:        request.CompanyName,
		CompanyAddress:     request.CompanyAddress,
		CompanyPhone:       request.CompanyPhone,
		CompanyEmail:       request.CompanyEmail,
		CompanyLogo:        request.CompanyLogo,
		CompanyWebsite:     request.CompanyWebsite,
		TaxID:              request.TaxID,
		FoundedDate:        request.FoundedDate,
		CompanyType:        request.CompanyType,
		CompanyDescription: request.CompanyDescription,
	}

	updatedCompanyProfile, err := s.companyProfileRepo.UpdateCompanyProfile(id, companyProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to update company profile: %v", err)
	}

	responseDTO, err := FormatResponseCompanyProfile(updatedCompanyProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to format company profile response: %v", err)
	}

	return responseDTO, nil
}

func (s *companyProfileService) DeleteCompanyProfile(id string) error {

	err := s.companyProfileRepo.DeleteCompanyProfile(id)
	if err != nil {
		return fmt.Errorf("failed to delete company profile: %v", err)
	}

	return nil
}
