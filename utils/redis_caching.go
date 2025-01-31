package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pahmiudahgede/senggoldong/config"
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
		return logAndReturnError("Error setting data in Redis", err)
	}

	log.Printf("Data stored in Redis with key: %s", key)
	return nil
}

func GetData(key string) (string, error) {
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {

		return "", nil
	} else if err != nil {

		return "", logAndReturnError("Error retrieving data from Redis", err)
	}
	return val, nil
}

func DeleteData(key string) error {
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		return logAndReturnError("Error deleting data from Redis", err)
	}
	log.Printf("Data deleted from Redis with key: %s", key)
	return nil
}

func CheckKeyExists(key string) (bool, error) {
	val, err := config.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, logAndReturnError("Error checking if key exists in Redis", err)
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

	return data, nil
}

func DeleteSessionData(userID string) error {
	sessionKey := "session:" + userID
	return DeleteData(sessionKey)
}

func logAndReturnError(message string, err error) error {
	log.Printf("%s: %v", message, err)
	return err
}
