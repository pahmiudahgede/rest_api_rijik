package main

import (
	"rijig/config"
	// "rijig/internal/repositories"
	// "rijig/internal/services"

	"rijig/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.SetupConfig()
	// cartRepo := repositories.NewCartRepository()
	// trashRepo := repositories.NewTrashRepository(config.DB)
	// cartService := services.NewCartService(cartRepo, trashRepo)
	// worker := worker.NewCartWorker(cartService, cartRepo, trashRepo)

	// go func() {
	// 	ticker := time.NewTicker(30 * time.Second)
	// 	defer ticker.Stop()

	// 	for range ticker.C {
	// 		if err := worker.AutoCommitExpiringCarts(); err != nil {
	// 			log.Printf("Auto-commit error: %v", err)
	// 		}
	// 	}
	// }()

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

	router.SetupRoutes(app)
	config.StartServer(app)
}
