package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	ctx := context.Background()
	redisClient := config.RedisClient

	redisKey := fmt.Sprintf("user_pin:%s", userID)

	pin, err := redisClient.Get(ctx, redisKey).Result()
	if err == nil {

		return domain.UserPin{
			UserID: userID,
			Pin:    pin,
		}, nil
	}

	var dbPin domain.UserPin
	err = config.DB.Where("user_id = ?", userID).First(&dbPin).Error
	if err != nil {
		return dbPin, errors.New("PIN tidak ditemukan")
	}

	redisClient.Set(ctx, redisKey, dbPin.Pin, 5*time.Minute)

	return dbPin, nil
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

	redisClient := config.RedisClient
	redisClient.Del(context.Background(), fmt.Sprintf("user_pin:%s", userID))

	return pin, nil
}
