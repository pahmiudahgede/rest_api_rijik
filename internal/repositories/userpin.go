package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"golang.org/x/crypto/bcrypt"
)

func CreatePin(pin *domain.UserPin) error {
	result := config.DB.Create(pin)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetPinByUserID(userID string) (domain.UserPin, error) {
	var pin domain.UserPin
	err := config.DB.Where("user_id = ?", userID).First(&pin).Error
	if err != nil {
		return pin, errors.New("PIN tidak ditemukan")
	}
	return pin, nil
}

func UpdatePin(userID string, newPin string) (domain.UserPin, error) {
	var pin domain.UserPin

	err := config.DB.Where("user_id = ?", userID).First(&pin).Error
	if err != nil {
		return pin, errors.New("PIN tidak ditemukan")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return pin, err
	}

	pin.Pin = string(hashedPin)

	if err := config.DB.Save(&pin).Error; err != nil {
		return pin, err
	}

	return pin, nil
}