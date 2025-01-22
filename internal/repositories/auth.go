package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func IsEmailExist(email, roleId string) bool {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("email:%s", email)
	cachedRole, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedRole == roleId {
		return true
	}

	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			if err := config.RedisClient.Set(ctx, cacheKey, roleId, 24*time.Hour).Err(); err != nil {
				log.Printf("Redis Set error: %v", err)
			}
			return true
		}
	}
	return false
}

func IsUsernameExist(username, roleId string) bool {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("username:%s", username)
	cachedRole, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedRole == roleId {
		return true
	}

	var user domain.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			if err := config.RedisClient.Set(ctx, cacheKey, roleId, 24*time.Hour).Err(); err != nil {
				log.Printf("Redis Set error: %v", err)
			}
			return true
		}
	}
	return false
}

func IsPhoneExist(phone, roleId string) bool {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("phone:%s", phone)
	cachedRole, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedRole == roleId {
		return true
	}

	var user domain.User
	if err := config.DB.Where("phone = ?", phone).First(&user).Error; err == nil {
		if user.RoleID == roleId {
			if err := config.RedisClient.Set(ctx, cacheKey, roleId, 24*time.Hour).Err(); err != nil {
				log.Printf("Redis Set error: %v", err)
			}
			return true
		}
	}
	return false
}

func CreateUser(username, name, email, phone, password, roleId string) error {
	user := domain.User{
		Username: username,
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
		RoleID:   roleId,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func GetUserByEmailUsernameOrPhone(identifier, roleId string) (domain.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", identifier)
	var user domain.User

	cachedUser, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			if roleId == "" || user.RoleID == roleId {
				return user, nil
			}
		}
	}

	err = config.DB.Where("email = ? OR username = ? OR phone = ?", identifier, identifier, identifier).First(&user).Error
	if err != nil {
		return user, errors.New("user not found")
	}

	if roleId != "" && user.RoleID != roleId {
		return user, errors.New("identifier found but role does not match")
	}

	userJSON, _ := json.Marshal(user)
	if err := config.RedisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour).Err(); err != nil {
		log.Printf("Redis Set error: %v", err)
	}
	return user, nil
}

func GetUserByID(userID string) (domain.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", userID)
	var user domain.User

	cachedUser, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return user, nil
		}
	}

	if err := config.DB.Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	userJSON, _ := json.Marshal(user)
	if err := config.RedisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour).Err(); err != nil {
		log.Printf("Redis Set error: %v", err)
	}

	return user, nil
}

func UpdateUser(user *domain.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return errors.New("failed to update user")
	}
	cacheKey := fmt.Sprintf("user:%s", user.ID)
	if err := config.RedisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		log.Printf("Redis Del error: %v", err)
	}
	return nil
}

func UpdateUserPassword(userID, newPassword string) error {
	var user domain.User

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.Password = newPassword

	if err := config.DB.Save(&user).Error; err != nil {
		return errors.New("failed to update password")
	}

	cacheKey := fmt.Sprintf("user:%s", userID)
	if err := config.RedisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		log.Printf("Redis Del error: %v", err)
	}

	return nil
}
