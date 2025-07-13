package userprofile

import (
	"context"
	"errors"
	"fmt"
	"rijig/internal/authentication"
	"rijig/internal/role"
	"rijig/model"
	"rijig/utils"
	"time"
)

var (
	ErrUserNotFound = errors.New("user tidak ditemukan")
)

type UserProfileService interface {
	GetUserProfile(ctx context.Context, userID string) (*UserProfileResponseDTO, error)
	UpdateRegistUserProfile(ctx context.Context, userID, deviceId string, req *RequestUserProfileDTO) (*authentication.AuthResponse, error)
}

type userProfileService struct {
	repo UserProfileRepository
}

func NewUserProfileService(repo UserProfileRepository) UserProfileService {
	return &userProfileService{
		repo: repo,
	}
}

func (s *userProfileService) GetUserProfile(ctx context.Context, userID string) (*UserProfileResponseDTO, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return s.mapToResponseDTO(user), nil
}

func (s *userProfileService) UpdateRegistUserProfile(ctx context.Context, userID, deviceId string, req *RequestUserProfileDTO) (*authentication.AuthResponse, error) {

	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	updateUser := &model.User{
		Name:                 req.Name,
		Gender:               req.Gender,
		Dateofbirth:          req.Dateofbirth,
		Placeofbirth:         req.Placeofbirth,
		Phone:                req.Phone,
		RegistrationProgress: utils.ProgressDataSubmitted,
	}

	if err := s.repo.Update(ctx, userID, updateUser); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	updatedUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	tokenResponse, err := utils.GenerateTokenPair(
		updatedUser.ID,
		updatedUser.Role.RoleName,
		deviceId,
		updatedUser.RegistrationStatus,
		int(updatedUser.RegistrationProgress),
	)

	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	nextStep := utils.GetNextRegistrationStep(
		updatedUser.Role.RoleName,
		int(updatedUser.RegistrationProgress),
		updateUser.RegistrationStatus,
	)

	return &authentication.AuthResponse{
		Message:            "Isi data diri berhasil",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: updateUser.RegistrationStatus,
		NextStep:           nextStep,
		SessionID:          tokenResponse.SessionID,
	}, nil
}

func (s *userProfileService) mapToResponseDTO(user *model.User) *UserProfileResponseDTO {

	createdAt, err := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	if err != nil {
		createdAt = user.CreatedAt.Format(time.RFC3339)
	}

	updatedAt, err := utils.FormatDateToIndonesianFormat(user.UpdatedAt)
	if err != nil {
		updatedAt = user.UpdatedAt.Format(time.RFC3339)
	}

	response := &UserProfileResponseDTO{
		ID:            user.ID,
		Name:          user.Name,
		Gender:        user.Gender,
		Dateofbirth:   user.Dateofbirth,
		Placeofbirth:  user.Placeofbirth,
		Phone:         user.Phone,
		PhoneVerified: user.PhoneVerified,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if user.Avatar != nil {
		response.Avatar = *user.Avatar
	}

	if user.Role != nil {
		roleCreatedAt, err := utils.FormatDateToIndonesianFormat(user.Role.CreatedAt)
		if err != nil {
			roleCreatedAt = user.Role.CreatedAt.Format(time.RFC3339)
		}

		roleUpdatedAt, err := utils.FormatDateToIndonesianFormat(user.Role.UpdatedAt)
		if err != nil {
			roleUpdatedAt = user.Role.UpdatedAt.Format(time.RFC3339)
		}

		response.Role = role.RoleResponseDTO{
			ID:        user.Role.ID,
			RoleName:  user.Role.RoleName,
			CreatedAt: roleCreatedAt,
			UpdatedAt: roleUpdatedAt,
		}
	}

	return response
}
