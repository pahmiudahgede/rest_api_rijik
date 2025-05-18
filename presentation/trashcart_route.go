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

	cartRepo := repositories.NewCartRepository()
	trashRepo := repositories.NewTrashRepository(config.DB)
	cartService := services.NewCartService(cartRepo, trashRepo)
	cartHandler := handler.NewCartHandler(cartService)

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware)
	cart.Post("/", cartHandler.CreateCart)
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/commit", cartHandler.CommitCart)
	cart.Delete("/:id", cartHandler.DeleteCart)
}
