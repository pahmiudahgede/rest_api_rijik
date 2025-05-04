package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {

	redisDBStr := os.Getenv("REDIS_DB")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to integer: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Redis connected successfully!")
}
