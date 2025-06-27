package userpin

import (
	"context"
	"errors"
	"fmt"
	"rijig/internal/authentication"
	"rijig/internal/userprofile"
	"rijig/model"
	"rijig/utils"

	"gorm.io/gorm"
)

type UserPinService interface {
	CreateUserPin(ctx context.Context, userID, deviceId string, dto *RequestPinDTO) (*authentication.AuthResponse, error)
	VerifyUserPin(ctx context.Context, userID, deviceID string, pin *RequestPinDTO) (*authentication.AuthResponse, error)
}

type userPinService struct {
	UserPinRepo     UserPinRepository
	authRepo        authentication.AuthenticationRepository
	userProfileRepo userprofile.UserProfileRepository
}

func NewUserPinService(UserPinRepo UserPinRepository,
	authRepo authentication.AuthenticationRepository,
	userProfileRepo userprofile.UserProfileRepository) UserPinService {
	return &userPinService{UserPinRepo, authRepo, userProfileRepo}
}

var (
	Pinhasbeencreated = "PIN already created"
)

func (s *userPinService) CreateUserPin(ctx context.Context, userID, deviceId string, dto *RequestPinDTO) (*authentication.AuthResponse, error) {

	_, err := s.UserPinRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("%v", Pinhasbeencreated)
	}

	hashed, err := utils.HashingPlainText(dto.Pin)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PIN: %w", err)
	}

	userPin := &model.UserPin{
		UserID: userID,
		Pin:    hashed,
	}

	if err := s.UserPinRepo.Create(ctx, userPin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to create pin: %w", err)
	}

	updates := map[string]interface{}{
		"registration_progress": utils.ProgressComplete,
		"registration_status":   utils.RegStatusComplete,
	}

	if err = s.authRepo.PatchUser(ctx, userID, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	updated, err := s.userProfileRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userprofile.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	tokenResponse, err := utils.GenerateTokenPair(
		updated.ID,
		updated.Role.RoleName,
		deviceId,
		updated.RegistrationStatus,
		int(updated.RegistrationProgress),
	)

	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	nextStep := utils.GetNextRegistrationStep(
		updated.Role.RoleName,
		int(updated.RegistrationProgress),
		updated.RegistrationStatus,
	)

	return &authentication.AuthResponse{
		Message:            "mantap semuanya completed",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: updated.RegistrationStatus,
		NextStep:           nextStep,
		SessionID:          tokenResponse.SessionID,
	}, nil
}

func (s *userPinService) VerifyUserPin(ctx context.Context, userID, deviceID string, pin *RequestPinDTO) (*authentication.AuthResponse, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	userPin, err := s.UserPinRepo.FindByUserID(ctx, userID)
	if err != nil || userPin == nil {
		return nil, fmt.Errorf("PIN not found")
	}

	if !utils.CompareHashAndPlainText(userPin.Pin, pin.Pin) {
		return nil, fmt.Errorf("PIN does not match, %s , %s", userPin.Pin, pin.Pin)
	}

	// roleName := strings.ToLower(user.Role.RoleName)

	// updated, err := s.userProfileRepo.GetByID(ctx, userID)
	// if err != nil {
	// 	if errors.Is(err, userprofile.ErrUserNotFound) {
	// 		return nil, fmt.Errorf("user not found")
	// 	}
	// 	return nil, fmt.Errorf("failed to get updated user: %w", err)
	// }

	tokenResponse, err := utils.GenerateTokenPair(
		user.ID,
		user.Role.RoleName,
		deviceID,
		user.RegistrationStatus,
		int(user.RegistrationProgress),
	)

	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	nextStep := utils.GetNextRegistrationStep(
		user.Role.RoleName,
		int(user.RegistrationProgress),
		user.RegistrationStatus,
	)

	return &authentication.AuthResponse{
		Message:            "mantap semuanya completed",
		AccessToken:        tokenResponse.AccessToken,
		RefreshToken:       tokenResponse.RefreshToken,
		TokenType:          string(tokenResponse.TokenType),
		ExpiresIn:          tokenResponse.ExpiresIn,
		RegistrationStatus: user.RegistrationStatus,
		NextStep:           nextStep,
		SessionID:          tokenResponse.SessionID,
	}, nil
}
