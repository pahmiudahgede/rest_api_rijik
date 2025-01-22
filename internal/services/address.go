package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
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

	cacheKey := fmt.Sprintf("address:user:%s", userID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return address, nil
}

func GetAllAddressesByUserID(userID string) ([]domain.Address, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("address:user:%s", userID)

	cachedAddresses, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var addresses []domain.Address
		if json.Unmarshal([]byte(cachedAddresses), &addresses) == nil {
			return addresses, nil
		}
	}

	addresses, err := repositories.GetAddressesByUserID(userID)
	if err != nil {
		return nil, err
	}

	addressesJSON, _ := json.Marshal(addresses)
	config.RedisClient.Set(ctx, cacheKey, addressesJSON, time.Hour).Err()

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

	cacheKey := fmt.Sprintf("address:user:%s", address.UserID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return updatedAddress, nil
}

func DeleteAddress(addressID string) error {
	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return errors.New("address not found")
	}

	err = repositories.DeleteAddress(addressID)
	if err != nil {
		return errors.New("failed to delete address")
	}

	cacheKey := fmt.Sprintf("address:user:%s", address.UserID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return nil
}
