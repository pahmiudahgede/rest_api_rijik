package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func CreateAddress(address *domain.Address) error {
	result := config.DB.Create(address)
	if result.Error != nil {
		return result.Error
	}

	cacheKey := fmt.Sprintf("address:user:%s", address.UserID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return nil
}

func GetAddressesByUserID(userID string) ([]domain.Address, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("address:user:%s", userID)

	cachedAddresses, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var addresses []domain.Address
		if json.Unmarshal([]byte(cachedAddresses), &addresses) == nil {
			return addresses, nil
		}
	}

	var addresses []domain.Address
	err = config.DB.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	addressesJSON, _ := json.Marshal(addresses)
	config.RedisClient.Set(ctx, cacheKey, addressesJSON, time.Hour).Err()

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

	cacheKey := fmt.Sprintf("address:user:%s", address.UserID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return address, nil
}

func DeleteAddress(addressID string) error {
	var address domain.Address
	if err := config.DB.Where("id = ?", addressID).First(&address).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&address).Error; err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("address:user:%s", address.UserID)
	config.RedisClient.Del(context.Background(), cacheKey)

	return nil
}
