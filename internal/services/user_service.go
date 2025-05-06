package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type UserService interface {
	GetUserByID(userID string) (*dto.UserResponseDTO, error)
	GetAllUsers(page, limit int) ([]dto.UserResponseDTO, error)
	UpdateUser(userID string, request *dto.RequestUserDTO) (*dto.UserResponseDTO, error)
	UpdateUserAvatar(userID string, avatar *multipart.FileHeader) (*dto.UserResponseDTO, error)
	UpdateUserPassword(userID, oldPassword, newPassword, confirmNewPassword string) error
}

type userService struct {
	userRepo repositories.UserProfilRepository
}

func NewUserService(userRepo repositories.UserProfilRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(userID string) (*dto.UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user by ID: %v", err)
	}

	userDTO, err := s.formatUserResponse(user)
	if err != nil {
		return nil, fmt.Errorf("error formatting user response: %v", err)
	}

	return userDTO, nil
}

func (s *userService) GetAllUsers(page, limit int) ([]dto.UserResponseDTO, error) {
	users, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all users: %v", err)
	}

	var userDTOs []dto.UserResponseDTO
	for _, user := range users {
		userDTO, err := s.formatUserResponse(&user)
		if err != nil {
			log.Printf("Error formatting user response for userID %s: %v", user.ID, err)
			continue
		}
		userDTOs = append(userDTOs, *userDTO)
	}

	return userDTOs, nil
}

func (s *userService) UpdateUser(userID string, request *dto.RequestUserDTO) (*dto.UserResponseDTO, error) {

	errors, valid := request.Validate()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	user.Name = request.Name
	user.Phone = request.Phone
	user.Email = request.Email

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	userDTO, err := s.formatUserResponse(user)
	if err != nil {
		return nil, fmt.Errorf("error formatting updated user response: %v", err)
	}

	return userDTO, nil
}

func (s *userService) UpdateUserAvatar(userID string, avatar *multipart.FileHeader) (*dto.UserResponseDTO, error) {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if *user.Avatar != "" {
		err := s.deleteAvatarImage(*user.Avatar)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old image: %v", err)
		}
	}

	avatarURL, err := s.saveAvatarImage(userID, avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar image: %v", err)
	}

	err = s.userRepo.UpdateAvatar(userID, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to update avatar in the database: %v", err)
	}

	userDTO, err := s.formatUserResponse(user)
	if err != nil {
		return nil, fmt.Errorf("failed to format user response: %v", err)
	}

	return userDTO, nil
}

func (s *userService) UpdateUserPassword(userID, oldPassword, newPassword, confirmNewPassword string) error {

	// errors, valid := utils.ValidatePasswordUpdate(oldPassword, newPassword, confirmNewPassword)
	// if !valid {
	// 	return fmt.Errorf("password validation error: %v", errors)
	// }

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.Password != oldPassword {
		return fmt.Errorf("old password is incorrect")
	}

	err = s.userRepo.UpdatePassword(userID, newPassword)
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}

	return nil
}

func (s *userService) formatUserResponse(user *model.User) (*dto.UserResponseDTO, error) {

	createdAt, _ := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	userDTO := &dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Name,
		Avatar:        user.Avatar,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.PhoneVerified,
		RoleName:      user.Role.RoleName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return userDTO, nil
}

func (s *userService) saveAvatarImage(userID string, avatar *multipart.FileHeader) (string, error) {

	pathImage := "/uploads/avatars/"
	avatarDir := "./public" + os.Getenv("BASE_URL") + pathImage

	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for avatar: %v", err)
		}
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(avatar.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	avatarFileName := fmt.Sprintf("%s_avatar%s", userID, extension)
	avatarPath := filepath.Join(avatarDir, avatarFileName)

	src, err := avatar.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(avatarPath)
	if err != nil {
		return "", fmt.Errorf("failed to create avatar file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save avatar: %v", err)
	}

	avatarURL := fmt.Sprintf("%s%s", pathImage, avatarFileName)

	return avatarURL, nil
}

func (s *userService) deleteAvatarImage(avatarPath string) error {

	if avatarPath == "" {
		return nil
	}

	baseDir := "./public/" + os.Getenv("BASE_URL")
	absolutePath := baseDir + avatarPath

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %v", err)
	}

	err := os.Remove(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to delete avatar image: %v", err)
	}

	log.Printf("Avatar image deleted successfully: %s", absolutePath)
	return nil
}
