package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

func GetPinByUserID(userID string) (domain.UserPin, error) {
	pin, err := repositories.GetPinByUserID(userID)
	if err != nil {
		return pin, errors.New("PIN tidak ditemukan")
	}
	return pin, nil
}

func CreatePin(userID string, input dto.PinInput) (domain.UserPin, error) {

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(input.Pin), bcrypt.DefaultCost)
	if err != nil {
		return domain.UserPin{}, err
	}

	pin := domain.UserPin{
		UserID: userID,
		Pin:    string(hashedPin),
	}

	err = repositories.CreatePin(&pin)
	if err != nil {
		return domain.UserPin{}, err
	}

	return pin, nil
}

func UpdatePin(userID string, oldPin string, newPin string) (domain.UserPin, error) {

	pin, err := repositories.GetPinByUserID(userID)
	if err != nil {
		return pin, errors.New("PIN tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pin.Pin), []byte(oldPin)); err != nil {
		return pin, errors.New("PIN lama tidak cocok")
	}

	updatedPin, err := repositories.UpdatePin(userID, newPin)
	if err != nil {
		return updatedPin, err
	}

	return updatedPin, nil
}

func CheckPin(storedPinHash string, inputPin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPinHash), []byte(inputPin))
	return err == nil
}
