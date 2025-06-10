package company

import (
	"context"
	"fmt"
	"log"
	"rijig/internal/authentication"
	"rijig/internal/role"
	"rijig/internal/userprofile"
	"rijig/model"
	"rijig/utils"
	"time"
)

type CompanyProfileService interface {
	CreateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error)
	GetCompanyProfileByID(ctx context.Context, id string) (*ResponseCompanyProfileDTO, error)
	GetCompanyProfilesByUserID(ctx context.Context, userID string) ([]ResponseCompanyProfileDTO, error)
	UpdateCompanyProfile(ctx context.Context, userID string, request *RequestCompanyProfileDTO) (*ResponseCompanyProfileDTO, error)
	DeleteCompanyProfile(ctx context.Context, userID string) error

	GetAllCompanyProfilesByRegStatus(ctx context.Context, userRegStatus string) ([]ResponseCompanyProfileDTO, error)
	UpdateUserRegistrationStatusByCompany(ctx context.Context, companyUserID string, newStatus string) error
}

type companyProfileService struct {
	companyRepo CompanyProfileRepository
	authRepo    authentication.AuthenticationRepository
}

func NewCompanyProfileService(companyRepo CompanyProfileRepository, authRepo authentication.AuthenticationRepository) CompanyProfileService {
	return &companyProfileService{
		companyRepo, authRepo,
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
	// if errors, valid := request.ValidateCompanyProfileInput(); !valid {
	// 	return nil, fmt.Errorf("validation failed: %v", errors)
	// }

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

func (s *companyProfileService) GetAllCompanyProfilesByRegStatus(ctx context.Context, userRegStatus string) ([]ResponseCompanyProfileDTO, error) {
	companyProfiles, err := s.authRepo.GetCompanyProfilesByUserRegStatus(ctx, userRegStatus)
	if err != nil {
		log.Printf("Error getting company profiles by registration status: %v", err)
		return nil, fmt.Errorf("failed to get company profiles: %w", err)
	}

	var response []ResponseCompanyProfileDTO
	for _, profile := range companyProfiles {
		dto := ResponseCompanyProfileDTO{
			ID:                 profile.ID,
			UserID:             profile.UserID,
			CompanyName:        profile.CompanyName,
			CompanyAddress:     profile.CompanyAddress,
			CompanyPhone:       profile.CompanyPhone,
			CompanyEmail:       profile.CompanyEmail,
			CompanyLogo:        profile.CompanyLogo,
			CompanyWebsite:     profile.CompanyWebsite,
			TaxID:              profile.TaxID,
			FoundedDate:        profile.FoundedDate,
			CompanyType:        profile.CompanyType,
			CompanyDescription: profile.CompanyDescription,
			CreatedAt:          profile.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          profile.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, dto)
	}

	return response, nil
}

func (s *companyProfileService) UpdateUserRegistrationStatusByCompany(ctx context.Context, companyUserID string, newStatus string) error {

	user, err := s.authRepo.FindUserByID(ctx, companyUserID)
	if err != nil {
		log.Printf("Error finding user by ID %s: %v", companyUserID, err)
		return fmt.Errorf("user not found: %w", err)
	}

	updates := map[string]interface{}{
		"registration_status": newStatus,
		"updated_at":          time.Now(),
	}

	switch newStatus {
	case utils.RegStatusConfirmed:
		updates["registration_progress"] = utils.ProgressDataSubmitted
	case utils.RegStatusRejected:
		updates["registration_progress"] = utils.ProgressOTPVerified
	}

	err = s.authRepo.PatchUser(ctx, user.ID, updates)
	if err != nil {
		log.Printf("Error updating user registration status for user ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to update user registration status: %w", err)
	}

	log.Printf("Successfully updated registration status for user ID %s to %s", user.ID, newStatus)
	return nil
}

func (s *companyProfileService) GetUserProfile(ctx context.Context, userID string) (*userprofile.UserProfileResponseDTO, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user profile for ID %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	response := &userprofile.UserProfileResponseDTO{
		ID:            user.ID,
		Name:          user.Name,
		Gender:        user.Gender,
		Dateofbirth:   user.Dateofbirth,
		Placeofbirth:  user.Placeofbirth,
		Phone:         user.Phone,
		Email:         user.Email,
		PhoneVerified: user.PhoneVerified,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.Avatar != nil {
		response.Avatar = *user.Avatar
	}

	if user.Role != nil {
		response.Role = role.RoleResponseDTO{
			ID:        user.Role.ID,
			RoleName:  user.Role.RoleName,
			CreatedAt: user.Role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.Role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response, nil
}
