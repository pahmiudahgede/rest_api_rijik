package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"rijig/config"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const defaultExpiration = 1 * time.Hour

func SetData[T any](key string, value T, expiration time.Duration) error {
	if expiration == 0 {
		expiration = defaultExpiration
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return logAndReturnError("Error marshaling data to JSON", err)
	}

	err = config.RedisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return logAndReturnError(fmt.Sprintf("Error setting data in Redis with key: %s", key), err)
	}

	log.Printf("Data stored in Redis with key: %s", key)
	return nil
}

func SaveSessionTokenToRedis(userID string, deviceID string, token string) error {

	sessionKey := "session:" + userID + ":" + deviceID

	err := config.RedisClient.Set(ctx, sessionKey, token, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	log.Printf("Session token saved to Redis with key: %s", sessionKey)
	return nil
}

func GetSessionTokenFromRedis(userID string, deviceID string) (string, error) {
	sessionKey := "session:" + userID + ":" + deviceID
	token, err := config.RedisClient.Get(ctx, sessionKey).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetData(key string) (string, error) {
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {

		return "", nil
	} else if err != nil {

		return "", logAndReturnError(fmt.Sprintf("Error retrieving data from Redis with key: %s", key), err)
	}
	return val, nil
}

func DeleteData(key string) error {
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		return logAndReturnError(fmt.Sprintf("Error deleting data from Redis with key: %s", key), err)
	}
	log.Printf("Data deleted from Redis with key: %s", key)
	return nil
}

func CheckKeyExists(key string) (bool, error) {
	val, err := config.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, logAndReturnError(fmt.Sprintf("Error checking if key exists in Redis with key: %s", key), err)
	}
	return val > 0, nil
}

func SetJSONData[T any](key string, value T, expiration time.Duration) error {
	return SetData(key, value, expiration)
}

func GetJSONData(key string) (map[string]interface{}, error) {
	val, err := GetData(key)
	if err != nil || val == "" {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Printf("Error unmarshaling JSON data from Redis: %v", err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("error: no data found for key %s", key)
	}

	return data, nil
}

func DeleteSessionData(userID string, deviceID string) error {
	sessionKey := "session:" + userID + ":" + deviceID
	sessionTokenKey := "session_token:" + userID + ":" + deviceID

	log.Printf("Attempting to delete session data with keys: %s, %s", sessionKey, sessionTokenKey)

	err := DeleteData(sessionKey)
	if err != nil {
		return fmt.Errorf("failed to delete session data: %w", err)
	}
	err = DeleteData(sessionTokenKey)
	if err != nil {
		return fmt.Errorf("failed to delete session token: %w", err)
	}

	log.Printf("Successfully deleted session data for userID: %s, deviceID: %s", userID, deviceID)
	return nil
}

func logAndReturnError(message string, err error) error {
	log.Printf("%s: %v", message, err)
	return err
}

func SetStringData(key, value string, expiration time.Duration) error {
	if expiration == 0 {
		expiration = defaultExpiration
	}

	err := config.RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return logAndReturnError(fmt.Sprintf("Error setting string data in Redis with key: %s", key), err)
	}

	log.Printf("String data stored in Redis with key: %s", key)
	return nil
}

func GetStringData(key string) (string, error) {
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", logAndReturnError(fmt.Sprintf("Error retrieving string data from Redis with key: %s", key), err)
	}

	return val, nil
}

func CheckSessionExists(userID string, deviceID string) (bool, error) {
	sessionKey := "session:" + userID + ":" + deviceID
	val, err := config.RedisClient.Exists(ctx, sessionKey).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
