package userpin

import (
	"context"
	"fmt"
	"rijig/internal/authentication"
	"rijig/model"
	"rijig/utils"
	"strings"
)

type UserPinService interface {
	CreateUserPin(ctx context.Context, userID string, dto *RequestPinDTO) error
	VerifyUserPin(ctx context.Context, userID string, pin *RequestPinDTO) (*utils.TokenResponse, error)
}

type userPinService struct {
	UserPinRepo UserPinRepository
	authRepo    authentication.AuthenticationRepository
}

func NewUserPinService(UserPinRepo UserPinRepository,
	authRepo authentication.AuthenticationRepository) UserPinService {
	return &userPinService{UserPinRepo, authRepo}
}

func (s *userPinService) CreateUserPin(ctx context.Context, userID string, dto *RequestPinDTO) error {

	if errs, ok := dto.ValidateRequestPinDTO(); !ok {
		return fmt.Errorf("validation error: %v", errs)
	}

	existingPin, err := s.UserPinRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing PIN: %w", err)
	}
	if existingPin != nil {
		return fmt.Errorf("PIN already created")
	}

	hashed, err := utils.HashingPlainText(dto.Pin)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	userPin := &model.UserPin{
		UserID: userID,
		Pin:    hashed,
	}

	if err := s.UserPinRepo.Create(ctx, userPin); err != nil {
		return fmt.Errorf("failed to create PIN: %w", err)
	}

	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	roleName := strings.ToLower(user.Role.RoleName)

	progress := authentication.IsRegistrationComplete(roleName, int(user.RegistrationProgress))
	// progress := utils.GetNextRegistrationStep(roleName, int(user.RegistrationProgress))
	// progress := utils.GetNextRegistrationStep(roleName, user.RegistrationProgress)
	// progress := utils.GetNextRegistrationStep(roleName, user.RegistrationProgress)

	if !progress {
		err = s.authRepo.PatchUser(ctx, userID, map[string]interface{}{
			"registration_progress": int(user.RegistrationProgress) + 1,
			"registration_status":   utils.RegStatusComplete,
		})
		if err != nil {
			return fmt.Errorf("failed to update user progress: %w", err)
		}
	}

	return nil
}

func (s *userPinService) VerifyUserPin(ctx context.Context, userID string, pin *RequestPinDTO) (*utils.TokenResponse, error) {
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

	roleName := strings.ToLower(user.Role.RoleName)
	return utils.GenerateTokenPair(user.ID, roleName, pin.DeviceId, user.RegistrationStatus, int(user.RegistrationProgress))
}
