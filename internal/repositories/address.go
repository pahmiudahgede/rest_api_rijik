package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func CreateAddress(address *domain.Address) error {
	result := config.DB.Create(address)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAddressesByUserID(userID string) ([]domain.Address, error) {
	var addresses []domain.Address
	err := config.DB.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func GetAddressByID(addressID string) (domain.Address, error) {
	var address domain.Address
	if err := config.DB.Where("id = ?", addressID).First(&address).Error; err != nil {
		return address, errors.New("address not found")
	}
	return address, nil
}

func UpdateAddress(address domain.Address) (domain.Address, error) {
	if err := config.DB.Save(&address).Error; err != nil {
		return address, err
	}
	return address, nil
}

func DeleteAddress(addressID string) error {
	var address domain.Address
	if err := config.DB.Where("id = ?", addressID).Delete(&address).Error; err != nil {
		return err
	}
	return nil
}
