package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func CreateAddress(userID string, input dto.AddressInput) (domain.Address, error) {
	address := domain.Address{
		UserID:      userID,
		Province:    input.Province,
		District:    input.District,
		Subdistrict: input.Subdistrict,
		PostalCode:  input.PostalCode,
		Village:     input.Village,
		Detail:      input.Detail,
		Geography:   input.Geography,
	}

	err := repositories.CreateAddress(&address)
	if err != nil {
		return domain.Address{}, err
	}

	return address, nil
}

func GetAllAddressesByUserID(userID string) ([]domain.Address, error) {

	addresses, err := repositories.GetAddressesByUserID(userID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func GetAddressByID(addressID string) (domain.Address, error) {
	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return address, errors.New("address not found")
	}
	return address, nil
}

func UpdateAddress(addressID string, input dto.AddressInput) (domain.Address, error) {

	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return address, errors.New("address not found")
	}

	address.Province = input.Province
	address.District = input.District
	address.Subdistrict = input.Subdistrict
	address.PostalCode = input.PostalCode
	address.Village = input.Village
	address.Detail = input.Detail
	address.Geography = input.Geography

	updatedAddress, err := repositories.UpdateAddress(address)
	if err != nil {
		return updatedAddress, errors.New("failed to update address")
	}

	return updatedAddress, nil
}

func DeleteAddress(addressID string) error {
	err := repositories.DeleteAddress(addressID)
	if err != nil {
		return errors.New("failed to delete address")
	}
	return nil
}
