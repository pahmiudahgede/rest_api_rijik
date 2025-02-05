package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserPinService interface {
	CreateUserPin(userID, pin string) (*dto.UserPinResponseDTO, error)
	VerifyUserPin(userID, pin string) (*dto.UserPinResponseDTO, error)
	CheckPinStatus(userID string) (string, *dto.UserPinResponseDTO, error)
	UpdateUserPin(userID, oldPin, newPin string) (*dto.UserPinResponseDTO, error)
}

type userPinService struct {
	UserPinRepo repositories.UserPinRepository
}

func NewUserPinService(userPinRepo repositories.UserPinRepository) UserPinService {
	return &userPinService{UserPinRepo: userPinRepo}
}

func (s *userPinService) VerifyUserPin(pin string, userID string) (*dto.UserPinResponseDTO, error) {

	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("user pin not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPin.Pin), []byte(pin))
	if err != nil {
		return nil, fmt.Errorf("incorrect pin")
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(userPin.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(userPin.UpdatedAt)

	userPinResponse := &dto.UserPinResponseDTO{
		ID:        userPin.ID,
		UserID:    userPin.UserID,
		Pin:       userPin.Pin,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return userPinResponse, nil
}

func (s *userPinService) CheckPinStatus(userID string) (string, *dto.UserPinResponseDTO, error) {
	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return "Pin not created", nil, nil
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(userPin.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(userPin.UpdatedAt)

	userPinResponse := &dto.UserPinResponseDTO{
		ID:        userPin.ID,
		UserID:    userPin.UserID,
		Pin:       userPin.Pin,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return "Pin already created", userPinResponse, nil
}

func (s *userPinService) CreateUserPin(userID, pin string) (*dto.UserPinResponseDTO, error) {

	existingPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil && existingPin != nil {
		return nil, fmt.Errorf("you have already created a pin, you don't need to create another one")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing the pin: %v", err)
	}

	newPin := model.UserPin{
		UserID: userID,
		Pin:    string(hashedPin),
	}

	err = s.UserPinRepo.Create(&newPin)
	if err != nil {
		return nil, fmt.Errorf("error creating user pin: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(newPin.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(newPin.UpdatedAt)

	userPinResponse := &dto.UserPinResponseDTO{
		ID:        newPin.ID,
		UserID:    newPin.UserID,
		Pin:       newPin.Pin,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	cacheKey := fmt.Sprintf("userpin:%s", userID)
	cacheData := map[string]interface{}{
		"data": userPinResponse,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching new user pin to Redis: %v\n", err)
	}

	return userPinResponse, nil
}

func (s *userPinService) UpdateUserPin(userID, oldPin, newPin string) (*dto.UserPinResponseDTO, error) {

	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("user pin not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPin.Pin), []byte(oldPin))
	if err != nil {
		return nil, fmt.Errorf("incorrect old pin")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing the new pin: %v", err)
	}

	userPin.Pin = string(hashedPin)
	err = s.UserPinRepo.Update(userPin)
	if err != nil {
		return nil, fmt.Errorf("error updating user pin: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(userPin.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(userPin.UpdatedAt)

	userPinResponse := &dto.UserPinResponseDTO{
		ID:        userPin.ID,
		UserID:    userPin.UserID,
		Pin:       userPin.Pin,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	cacheKey := fmt.Sprintf("userpin:%s", userID)
	cacheData := map[string]interface{}{
		"data": userPinResponse,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated user pin to Redis: %v\n", err)
	}

	return userPinResponse, nil
}
