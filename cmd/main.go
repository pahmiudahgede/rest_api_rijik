package main

import (
	"context"
	"log"
	"rijig/config"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/router"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.SetupConfig()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,x-api-key",
	}))

	// Route setup
	router.SetupRoutes(app)

	// Siapkan dependency untuk worker
	repoCart := repositories.NewCartRepository()
	repoTrash := repositories.NewTrashRepository(config.DB)
	cartService := services.NewCartService(repoCart, repoTrash)
	ctx := context.Background()

	// ✅ Jalankan worker di background
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		log.Println("🛠️ Cart Worker is running in background...")
		for range ticker.C {
			processCartKeys(ctx, cartService)
		}
	}()

	// 🚀 Jalankan server (blocking)
	config.StartServer(app)
}

func processCartKeys(ctx context.Context, cartService services.CartService) {
	pattern := "cart:user:*"
	iter := config.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := config.RedisClient.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get TTL for key %s: %v", key, err)
			continue
		}

		if ttl <= time.Minute {
			log.Printf("🔄 Auto-committing key: %s", key)
			parts := strings.Split(key, ":")
			if len(parts) != 3 {
				log.Printf("Invalid key format: %s", key)
				continue
			}
			userID := parts[2]

			err := cartService.CommitCartFromRedis(userID)
			if err != nil {
				log.Printf("❌ Failed to commit cart for user %s: %v", userID, err)
			} else {
				log.Printf("✅ Cart for user %s committed successfully", userID)
			}
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Error iterating keys: %v", err)
	}
}
