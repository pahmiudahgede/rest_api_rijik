package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type UserProfileService interface {
	GetUserProfile(userID string) (*dto.UserResponseDTO, error)
}

type userProfileService struct {
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
