package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
)

var allowedExtensions = []string{".jpg", ".jpeg", ".png"}

type UserProfileService interface {
	GetUserProfile(userID string) (*dto.UserResponseDTO, error)
	UpdateUserProfile(userID string, updateData dto.UpdateUserDTO) (*dto.UserResponseDTO, error)
	UpdateUserPassword(userID string, passwordData dto.UpdatePasswordDTO) (string, error)
	UpdateUserAvatar(userID string, file *multipart.FileHeader) (string, error)

	GetAllUsers() ([]dto.UserResponseDTO, error)
	GetUsersByRoleID(roleID string) ([]dto.UserResponseDTO, error)
}

type userProfileService struct {
	UserRepo        repositories.UserRepository
	RoleRepo        repositories.RoleRepository
	UserProfileRepo repositories.UserProfileRepository
}

func NewUserProfileService(userProfileRepo repositories.UserProfileRepository) UserProfileService {
	return &userProfileService{UserProfileRepo: userProfileRepo}
}

func (s *userProfileService) prepareUserResponse(user *model.User) *dto.UserResponseDTO {
	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	return &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Avatar:        user.Avatar,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		RoleName:      user.Role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
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

	userResponse := s.prepareUserResponse(user)

	cacheData := map[string]interface{}{
		"data": userResponse,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching user profile to Redis: %v\n", err)
	}

	return userResponse, nil
}

func (s *userProfileService) GetAllUsers() ([]dto.UserResponseDTO, error) {
	users, err := s.UserProfileRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponseDTO
	for _, user := range users {
		response = append(response, dto.UserResponseDTO{
			ID:            user.ID,
			Username:      user.Username,
			Avatar:        user.Avatar,
			Name:          user.Name,
			Phone:         user.Phone,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			RoleName:      user.Role.RoleName,
			CreatedAt:     user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *userProfileService) GetUsersByRoleID(roleID string) ([]dto.UserResponseDTO, error) {
	users, err := s.UserProfileRepo.FindByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponseDTO
	for _, user := range users {
		response = append(response, dto.UserResponseDTO{
			ID:            user.ID,
			Username:      user.Username,
			Avatar:        user.Avatar,
			Name:          user.Name,
			Phone:         user.Phone,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			RoleName:      user.Role.RoleName,
			CreatedAt:     user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
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
		if err := s.updatePhoneIfNeeded(user, updateData.Phone); err != nil {
			return nil, err
		}
		user.Phone = updateData.Phone
	}

	if updateData.Email != "" && updateData.Email != user.Email {
		if err := s.updateEmailIfNeeded(user, updateData.Email); err != nil {
			return nil, err
		}
		user.Email = updateData.Email
	}

	err = s.UserProfileRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	userResponse := s.prepareUserResponse(user)

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

func (s *userProfileService) updatePhoneIfNeeded(user *model.User, newPhone string) error {
	existingPhone, _ := s.UserRepo.FindByPhoneAndRole(newPhone, user.RoleID)
	if existingPhone != nil {
		return fmt.Errorf("phone number is already used for this role")
	}
	return nil
}

func (s *userProfileService) updateEmailIfNeeded(user *model.User, newEmail string) error {
	existingEmail, _ := s.UserRepo.FindByEmailAndRole(newEmail, user.RoleID)
	if existingEmail != nil {
		return fmt.Errorf("email is already used for this role")
	}
	return nil
}

func (s *userProfileService) UpdateUserPassword(userID string, passwordData dto.UpdatePasswordDTO) (string, error) {

	validationErrors, valid := passwordData.Validate()
	if !valid {
		return "", fmt.Errorf("validation failed: %v", validationErrors)
	}

	user, err := s.UserProfileRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !CheckPasswordHash(passwordData.OldPassword, user.Password) {
		return "", errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash new password: %v", err)
	}

	user.Password = string(hashedPassword)
	err = s.UserProfileRepo.Update(user)
	if err != nil {
		return "", fmt.Errorf("failed to update password: %v", err)
	}

	return "Password berhasil diupdate", nil
}

func (s *userProfileService) UpdateUserAvatar(userID string, file *multipart.FileHeader) (string, error) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return "", fmt.Errorf("BASE_URL is not set in environment variables")
	}

	avatarDir := filepath.Join("./public", baseURL, "/uploads/avatars")
	if err := ensureAvatarDirectoryExists(avatarDir); err != nil {
		return "", err
	}

	if err := validateAvatarFile(file); err != nil {
		return "", err
	}

	updatedUser, err := s.UserProfileRepo.FindByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user data: %v", err)
	}

	if updatedUser.Avatar != nil && *updatedUser.Avatar != "" {
		oldAvatarPath := filepath.Join("./public", *updatedUser.Avatar)
		if _, err := os.Stat(oldAvatarPath); err == nil {

			if err := os.Remove(oldAvatarPath); err != nil {
				return "", fmt.Errorf("failed to remove old avatar: %v", err)
			}
		} else {

			log.Printf("Old avatar file not found: %s", oldAvatarPath)
		}
	}

	avatarURL, err := saveAvatarFile(file, userID, avatarDir)
	if err != nil {
		return "", err
	}

	err = s.UserProfileRepo.UpdateAvatar(userID, avatarURL)
	if err != nil {
		return "", fmt.Errorf("failed to update avatar in the database: %v", err)
	}

	return "Foto profil berhasil diupdate", nil
}

func ensureAvatarDirectoryExists(avatarDir string) error {
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create avatar directory: %v", err)
		}
	}
	return nil
}

func validateAvatarFile(file *multipart.FileHeader) error {
	extension := filepath.Ext(file.Filename)
	for _, ext := range allowedExtensions {
		if extension == ext {
			return nil
		}
	}
	return fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
}

func saveAvatarFile(file *multipart.FileHeader, userID, avatarDir string) (string, error) {
	extension := filepath.Ext(file.Filename)
	avatarFileName := fmt.Sprintf("%s_avatar%s", userID, extension)
	avatarPath := filepath.Join(avatarDir, avatarFileName)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(avatarPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", fmt.Errorf("failed to save avatar file: %v", err)
	}

	relativePath := filepath.Join("/uploads/avatars", avatarFileName)
	return relativePath, nil
}
