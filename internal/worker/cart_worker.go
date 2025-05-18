package worker

import (
	"log"
	"strings"
	"time"

	"rijig/config"
	"rijig/internal/services"
)

const (
	lockPrefix      = "lock:cart:"
	lockExpiration  = 30 * time.Second
	commitThreshold = 20 * time.Second
	scanPattern     = "cart:*"
)

func StartCartCommitWorker(service *services.CartService) {

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		log.Println("🛠️ Cart Worker is running in background...")
		for range ticker.C {
			processCarts(service)
		}
	}()

}

func processCarts(service *services.CartService) {
	iter := config.RedisClient.Scan(config.Ctx, 0, scanPattern, 0).Iterator()
	for iter.Next(config.Ctx) {
		key := iter.Val()

		ttl, err := config.RedisClient.TTL(config.Ctx, key).Result()
		if err != nil {
			log.Printf("❌ Error getting TTL for %s: %v", key, err)
			continue
		}

		if ttl > 0 && ttl < commitThreshold {
			userID := extractUserIDFromKey(key)
			if userID == "" {
				continue
			}

			lockKey := lockPrefix + userID
			acquired, err := config.RedisClient.SetNX(config.Ctx, lockKey, "locked", lockExpiration).Result()
			if err != nil || !acquired {
				continue
			}

			log.Printf("🔄 Auto-committing cart for user %s (TTL: %v)", userID, ttl)
			if err := service.CommitCartToDatabase(userID); err != nil {
				log.Printf("❌ Failed to commit cart for %s: %v", userID, err)
			} else {
				log.Printf("✅ Cart committed for user %s", userID)
			}

		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("❌ Error iterating Redis keys: %v", err)
	}
}

func extractUserIDFromKey(key string) string {
	if strings.HasPrefix(key, "cart:") {
		return strings.TrimPrefix(key, "cart:")
	}
	return ""
}
