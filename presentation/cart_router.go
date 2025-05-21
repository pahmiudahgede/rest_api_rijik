package presentation

import (
	"rijig/internal/handler"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func TrashCartRouter(api fiber.Router) {
	cartService := services.NewCartService()
	cartHandler := handler.NewCartHandler(cartService)

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware)

	cart.Get("/", cartHandler.GetCart)
	cart.Post("/item", cartHandler.AddOrUpdateCartItem)
	cart.Post("/items", cartHandler.AddMultipleCartItems)
	cart.Delete("/item/:trashID", cartHandler.DeleteCartItem)
	cart.Delete("/", cartHandler.ClearCart)

}
