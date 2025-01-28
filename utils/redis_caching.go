package utils

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pahmiudahgede/senggoldong/config"
)

func SetData(key string, value interface{}, expiration time.Duration) error {

	err := config.RedisClient.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting data to Redis: %v", err)
		return err
	}
	log.Printf("Data stored in Redis with key: %s", key)
	return nil
}

func GetData(key string) (string, error) {

	val, err := config.RedisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		log.Printf("No data found for key: %s", key)
		return "", nil
	} else if err != nil {
		log.Printf("Error getting data from Redis: %v", err)
		return "", err
	}
	log.Printf("Data retrieved from Redis for key: %s", key)
	return val, nil
}

func DeleteData(key string) error {

	err := config.RedisClient.Del(context.Background(), key).Err()
	if err != nil {
		log.Printf("Error deleting data from Redis: %v", err)
		return err
	}
	log.Printf("Data deleted from Redis with key: %s", key)
	return nil
}

func SetDataWithExpire(key string, value interface{}, expiration time.Duration) error {

	err := config.RedisClient.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting data with expiration to Redis: %v", err)
		return err
	}
	log.Printf("Data stored in Redis with key: %s and expiration: %v", key, expiration)
	return nil
}

func CheckKeyExists(key string) (bool, error) {

	val, err := config.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		log.Printf("Error checking if key exists in Redis: %v", err)
		return false, err
	}
	return val > 0, nil
}
