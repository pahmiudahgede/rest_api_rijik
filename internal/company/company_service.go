package company

import (
	"context"
	"fmt"
	"rijig/model"
	"rijig/utils"
)

type CompanyProfileService interface {
	CreateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error)
	GetCompanyProfileByID(ctx context.Context, id string) (*ResponseCompanyProfileDTO, error)
	GetCompanyProfilesByUserID(ctx context.Context, userID string) ([]ResponseCompanyProfileDTO, error)
	UpdateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error)
	DeleteCompanyProfile(ctx context.Context, userID string) error
}

type companyProfileService struct {
	companyRepo CompanyProfileRepository
}

func NewCompanyProfileService(companyRepo CompanyProfileRepository) CompanyProfileService {
	return &companyProfileService{
		companyRepo: companyRepo,
	}
}

func FormatResponseCompanyProfile(companyProfile *model.CompanyProfile) (*ResponseCompanyProfileDTO, error) {
	createdAt, _ := utils.FormatDateToIndonesianFormat(companyProfile.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(companyProfile.UpdatedAt)

	return &ResponseCompanyProfileDTO{
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
	}, nil
}

func (s *companyProfileService) CreateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error) {
	if errors, valid := request.ValidateCompanyProfileInput(); !valid {
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

	created, err := s.companyRepo.CreateCompanyProfile(ctx, companyProfile)
	if err != nil {
		return nil, err
	}

	return FormatResponseCompanyProfile(created)
}

func (s *companyProfileService) GetCompanyProfileByID(ctx context.Context, id string) (*ResponseCompanyProfileDTO, error) {
	profile, err := s.companyRepo.GetCompanyProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return FormatResponseCompanyProfile(profile)
}

func (s *companyProfileService) GetCompanyProfilesByUserID(ctx context.Context, userID string) ([]ResponseCompanyProfileDTO, error) {
	profiles, err := s.companyRepo.GetCompanyProfilesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []ResponseCompanyProfileDTO
	for _, p := range profiles {
		dto, err := FormatResponseCompanyProfile(&p)
		if err != nil {
			continue
		}
		responses = append(responses, *dto)
	}

	return responses, nil
}

func (s *companyProfileService) UpdateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error) {
	if errors, valid := request.ValidateCompanyProfileInput(); !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	company := &model.CompanyProfile{
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

	if err := s.companyRepo.UpdateCompanyProfile(ctx, company); err != nil {
		return nil, err
	}

	updated, err := s.companyRepo.GetCompanyProfilesByUserID(ctx, userID)
	if err != nil || len(updated) == 0 {
		return nil, fmt.Errorf("failed to retrieve updated company profile")
	}

	return FormatResponseCompanyProfile(&updated[0])
}

func (s *companyProfileService) DeleteCompanyProfile(ctx context.Context, userID string) error {
	return s.companyRepo.DeleteCompanyProfileByUserID(ctx, userID)
}
