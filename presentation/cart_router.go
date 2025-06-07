package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func TrashCartRouter(api fiber.Router) {
	repo := repositories.NewCartRepository()
	trashRepo := repositories.NewTrashRepository(config.DB)
	cartService := services.NewCartService(repo, trashRepo)
	cartHandler := handler.NewCartHandler(cartService)

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware())

	cart.Get("/", cartHandler.GetCart)
	cart.Post("/item", cartHandler.AddOrUpdateItem)
	cart.Delete("/item/:trash_id", cartHandler.DeleteItem)
	cart.Delete("/clear", cartHandler.ClearCart)
}

// cart.Post("/items", cartHandler.AddMultipleCartItems)
