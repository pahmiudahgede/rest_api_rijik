package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserPinService interface {
	CreateUserPin(userID, pin string) (string, error)
	VerifyUserPin(userID, pin string) (string, error)
	CheckPinStatus(userID string) (string, error)
	UpdateUserPin(userID, oldPin, newPin string) (string, error)
}

type userPinService struct {
	UserPinRepo repositories.UserPinRepository
}

func NewUserPinService(userPinRepo repositories.UserPinRepository) UserPinService {
	return &userPinService{UserPinRepo: userPinRepo}
}

func (s *userPinService) VerifyUserPin(userID, pin string) (string, error) {
	if pin == "" {
		return "", fmt.Errorf("pin tidak boleh kosong")
	}

	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error fetching user pin: %v", err)
	}
	if userPin == nil {
		return "", fmt.Errorf("user pin not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPin.Pin), []byte(pin))
	if err != nil {
		return "", fmt.Errorf("incorrect pin")
	}

	return "Pin yang anda masukkan benar", nil
}

func (s *userPinService) CheckPinStatus(userID string) (string, error) {
	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error checking pin status: %v", err)
	}
	if userPin == nil {
		return "Pin not created", nil
	}

	return "Pin already created", nil
}

func (s *userPinService) CreateUserPin(userID, pin string) (string, error) {

	existingPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error checking existing pin: %v", err)
	}

	if existingPin != nil {
		return "", fmt.Errorf("you have already created a pin, you don't need to create another one")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing the pin: %v", err)
	}

	newPin := model.UserPin{
		UserID: userID,
		Pin:    string(hashedPin),
	}

	err = s.UserPinRepo.Create(&newPin)
	if err != nil {
		return "", fmt.Errorf("error creating user pin: %v", err)
	}

	cacheKey := fmt.Sprintf("userpin:%s", userID)
	cacheData := map[string]interface{}{"data": newPin}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching new user pin to Redis: %v\n", err)
	}

	return "Pin berhasil dibuat", nil
}

func (s *userPinService) UpdateUserPin(userID, oldPin, newPin string) (string, error) {

	userPin, err := s.UserPinRepo.FindByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error fetching user pin: %v", err)
	}

	if userPin == nil {
		return "", fmt.Errorf("user pin not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPin.Pin), []byte(oldPin))
	if err != nil {
		return "", fmt.Errorf("incorrect old pin")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing the new pin: %v", err)
	}

	userPin.Pin = string(hashedPin)
	err = s.UserPinRepo.Update(userPin)
	if err != nil {
		return "", fmt.Errorf("error updating user pin: %v", err)
	}

	cacheKey := fmt.Sprintf("userpin:%s", userID)
	cacheData := map[string]interface{}{"data": userPin}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated user pin to Redis: %v\n", err)
	}

	return "Pin berhasil diperbarui", nil
}
