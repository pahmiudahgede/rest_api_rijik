package main

import (
	"log"
	"rijig/config"
	"rijig/internal/cart"
	"rijig/internal/trash"
	"rijig/internal/worker"
	"time"

	// "rijig/internal/repositories"
	// "rijig/internal/services"

	"rijig/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.SetupConfig()
	cartRepo := cart.NewCartRepository()
	trashRepo := trash.NewTrashRepository(config.DB)
	cartService := cart.NewCartService(cartRepo, trashRepo)
	worker := worker.NewCartWorker(cartService, cartRepo, trashRepo)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := worker.AutoCommitExpiringCarts(); err != nil {
				log.Printf("Auto-commit error: %v", err)
			}
		}
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New())

	app.Get("/apirijig/v2/health", func(c *fiber.Ctx) error {
		// Check database connection
		db, err := config.DB.DB()
		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":  "unhealthy",
				"error":   "database connection failed",
				"service": "rijig-api",
				"version": "v2.0.0",
				"time":    time.Now(),
			})
		}

		if err := db.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":  "unhealthy",
				"error":   "database ping failed",
				"service": "rijig-api",
				"version": "v2.0.0",
				"time":    time.Now(),
			})
		}

		// Check Redis connection
		if _, err := config.RedisClient.Ping(config.Ctx).Result(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":  "unhealthy",
				"error":   "redis connection failed",
				"service": "rijig-api",
				"version": "v2.0.0",
				"time":    time.Now(),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "healthy",
			"service":  "rijig-api",
			"version":  "v2.0.0",
			"time":     time.Now(),
			"database": "connected",
			"redis":    "connected",
		})
	})

	// Simple ping endpoint
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
			"time":    time.Now(),
		})
	})

	router.SetupRoutes(app)
	config.StartServer(app)
}
