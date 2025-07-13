package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	// Check if running on Railway (REDIS_URL provided)
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		log.Println("Using Railway REDIS_URL")
		connectRedisFromURL(redisURL)
	} else {
		// Fallback to individual environment variables (for local development)
		log.Println("Using individual Redis environment variables")
		connectRedisFromEnvVars()
	}

	// Test connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Redis connected successfully!")
}

func connectRedisFromURL(redisURL string) {
	// Parse Redis URL
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("Error parsing REDIS_URL: %v", err)
	}

	// Extract password
	password := ""
	if parsedURL.User != nil {
		password, _ = parsedURL.User.Password()
	}

	// Extract database number
	db := 0
	if parsedURL.Path != "" && len(parsedURL.Path) > 1 {
		if dbNum, err := strconv.Atoi(parsedURL.Path[1:]); err == nil {
			db = dbNum
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     parsedURL.Host,
		Password: password,
		DB:       db,
	})
}

func connectRedisFromEnvVars() {
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Warning: Error converting REDIS_DB to integer, using default 0: %v", err)
		redisDB = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
}