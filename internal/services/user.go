package services

import (
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetUsers() ([]dto.UserResponseDTO, error) {
	users, err := repositories.GetUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponseDTO
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponseDTO{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			RoleId:    user.RoleID,
			CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
			UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
		})
	}
	return userResponses, nil
}

func GetUsersByRole(roleID string) ([]dto.UserResponseDTO, error) {
	users, err := repositories.GetUsersByRole(roleID)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponseDTO
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponseDTO{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			RoleId:    user.RoleID,
			CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
			UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
		})
	}
	return userResponses, nil
}

func GetUserByUserID(userID string) (dto.UserResponseDTO, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userResponse := dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		RoleId:    user.RoleID,
		CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
		UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
	}

	return userResponse, nil
}