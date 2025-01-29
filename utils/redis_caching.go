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

func SetData(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return err
	}

	err = config.RedisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		log.Printf("Error setting JSON data to Redis: %v", err)
		return err
	}

	log.Printf("JSON Data stored in Redis with key: %s", key)
	return nil
}

func GetData(ctx context.Context, key string) (string, error) {
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		log.Printf("Error getting data from Redis: %v", err)
		return "", err
	}
	return val, nil
}

func DeleteData(key string) error {
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting data from Redis: %v", err)
		return err
	}
	log.Printf("Data deleted from Redis with key: %s", key)
	return nil
}

func CheckKeyExists(ctx context.Context, key string) (bool, error) {
	val, err := config.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if key exists in Redis: %v", err)
		return false, err
	}
	return val > 0, nil
}

func SetJSONData(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return err
	}

	err = config.RedisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		log.Printf("Error setting JSON data to Redis: %v", err)
		return err
	}

	log.Printf("JSON Data stored in Redis with key: %s", key)
	return nil
}

func GetJSONData(ctx context.Context, key string) (map[string]interface{}, error) {
	val, err := GetData(ctx, key)
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
