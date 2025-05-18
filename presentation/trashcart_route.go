package presentation

import (
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/internal/worker"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func TrashCartRouter(api fiber.Router) {
	cartRepo := repositories.NewCartRepository()
	cartService := services.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	worker.StartCartCommitWorker(cartService)

	cart := api.Group("/cart", middleware.AuthMiddleware)
	cart.Put("/refresh", cartHandler.RefreshCartTTL)
	cart.Post("/", cartHandler.AddOrUpdateCartItem)
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/commit", cartHandler.CommitCart)
	cart.Delete("/", cartHandler.ClearCart)
	cart.Delete("/:trashid", cartHandler.DeleteCartItem)
}
