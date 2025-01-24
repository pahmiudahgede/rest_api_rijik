package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func Context() context.Context {
	return context.Background()
}

func InitRedis() {

	InitConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
		Password: RedisPassword,
		DB:       RedisDB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error: Gagal terhubung ke Redis: %v", err)
	}

	log.Println("Koneksi ke Redis berhasil.")
}

func GetFromCache(key string) (string, error) {
	val, err := RedisClient.Get(Context(), key).Result()
	if err == redis.Nil {

		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func SetToCache(key string, value string, ttl time.Duration) error {
	err := RedisClient.Set(Context(), key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteFromCache(key string) error {
	err := RedisClient.Del(Context(), key).Err()
	if err != nil {
		return err
	}
	return nil
}
