package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserProfileService interface {
	GetUserProfile(userID string) (*dto.UserResponseDTO, error)
	UpdateUserProfile(userID string, updateData dto.UpdateUserDTO) (*dto.UserResponseDTO, error)
	UpdateUserPassword(userID string, passwordData dto.UpdatePasswordDTO) (*dto.UserResponseDTO, error)
}

type userProfileService struct {
	UserRepo        repositories.UserRepository
	RoleRepo        repositories.RoleRepository
	UserProfileRepo repositories.UserProfileRepository
}

func NewUserProfileService(userProfileRepo repositories.UserProfileRepository) UserProfileService {
	return &userProfileService{UserProfileRepo: userProfileRepo}
}

func (s *userProfileService) GetUserProfile(userID string) (*dto.UserResponseDTO, error) {

	cacheKey := fmt.Sprintf("userProfile:%s", userID)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {

		userResponse := &dto.UserResponseDTO{}

		if data, ok := cachedData["data"].(string); ok {

			if err := json.Unmarshal([]byte(data), userResponse); err != nil {
				return nil, err
			}
			return userResponse, nil
		}
	}

	user, err := s.UserProfileRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	userResponse := &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		RoleName:      user.Role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": userResponse,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {

		fmt.Printf("Error caching user profile to Redis: %v\n", err)
	}

	return userResponse, nil
}

func (s *userProfileService) UpdateUserProfile(userID string, updateData dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {

	user, err := s.UserProfileRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	validationErrors, valid := updateData.Validate()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	if updateData.Name != "" {
		user.Name = updateData.Name
	}

	if updateData.Phone != "" && updateData.Phone != user.Phone {

		existingPhone, _ := s.UserRepo.FindByPhoneAndRole(updateData.Phone, user.RoleID)
		if existingPhone != nil {
			return nil, fmt.Errorf("phone number is already used for this role")
		}
		user.Phone = updateData.Phone
	}

	if updateData.Email != "" && updateData.Email != user.Email {

		existingEmail, _ := s.UserRepo.FindByEmailAndRole(updateData.Email, user.RoleID)
		if existingEmail != nil {
			return nil, fmt.Errorf("email is already used for this role")
		}
		user.Email = updateData.Email
	}

	err = s.UserProfileRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	userResponse := &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		RoleName:      user.Role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	cacheKey := fmt.Sprintf("userProfile:%s", userID)
	cacheData := map[string]interface{}{
		"data": userResponse,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error updating cached user profile in Redis: %v\n", err)
	}

	return userResponse, nil
}

func (s *userProfileService) UpdateUserPassword(userID string, passwordData dto.UpdatePasswordDTO) (*dto.UserResponseDTO, error) {

	validationErrors, valid := passwordData.Validate()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	user, err := s.UserProfileRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !CheckPasswordHash(passwordData.OldPassword, user.Password) {
		return nil, errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new password: %v", err)
	}

	user.Password = string(hashedPassword)

	err = s.UserProfileRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update password: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	userResponse := &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		RoleName:      user.Role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return userResponse, nil
}
