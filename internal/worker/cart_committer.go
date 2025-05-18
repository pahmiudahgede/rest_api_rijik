package worker

import (
	"context"
	"encoding/json"
	"log"

	"rijig/config"
	"rijig/internal/repositories"
	"rijig/model"
)

type CartCommitter struct {
	repo repositories.CartRepository
}

func NewCartCommitter(repo repositories.CartRepository) *CartCommitter {
	return &CartCommitter{repo: repo}
}

func (cc *CartCommitter) RunAutoCommit() {
	ctx := context.Background()
	pattern := "cart:user:*"

	iter := config.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		val, err := config.RedisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Error fetching key %s: %v", key, err)
			continue
		}

		var cart model.Cart
		if err := json.Unmarshal([]byte(val), &cart); err != nil {
			log.Printf("Invalid cart format in key %s: %v", key, err)
			continue
		}

		// Simpan ke DB
		if err := cc.repo.Create(&cart); err != nil {
			log.Printf("Failed to commit cart to DB from key %s: %v", key, err)
			continue
		}

		// Delete from Redis
		if err := config.RedisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("Failed to delete key %s after commit: %v", key, err)
		} else {
			log.Printf("Committed and deleted key %s successfully", key)
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Redis scan error: %v", err)
	}
}
