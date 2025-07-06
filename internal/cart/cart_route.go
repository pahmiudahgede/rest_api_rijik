package cart

import (
	"rijig/config"
	"rijig/internal/trash"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func TrashCartRouter(api fiber.Router) {
	repo := NewCartRepository()
	trashRepo := trash.NewTrashRepository(config.DB)
	cartService := NewCartService(repo, trashRepo)
	cartHandler := NewCartHandler(cartService)

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware())

	cart.Get("/", cartHandler.GetCart)
	cart.Post("/item", cartHandler.AddOrUpdateItem)
	cart.Delete("/item/:trash_id", cartHandler.DeleteItem)
	cart.Delete("/clear", cartHandler.ClearCart)
}
