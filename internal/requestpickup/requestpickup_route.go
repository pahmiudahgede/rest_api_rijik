package requestpickup

import (
	"rijig/config"
	"rijig/internal/cart"
	"rijig/internal/trash"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func RequestPickupRouter(api fiber.Router) {
	cartRepo := cart.NewCartRepository()
	pickupRepo := NewRequestPickupRepository()
	historyRepo := NewPickupStatusHistoryRepository()
	trashRepo := trash.NewTrashRepository(config.DB)

	cartService := cart.NewCartService(cartRepo, trashRepo)
	historyService := NewPickupStatusHistoryService(historyRepo)

	pickupService := NewRequestPickupService(trashRepo, pickupRepo, cartService, historyRepo)
	pickupHandler := NewRequestPickupHandler(pickupService)
	statuspickupHandler := NewPickupStatusHistoryHandler(historyService)

	reqpickup := api.Group("/reqpickup")
	reqpickup.Use(middleware.AuthMiddleware())

	reqpickup.Post("/manual", pickupHandler.CreateRequestPickup)
	reqpickup.Get("/pickup/:id/history", statuspickupHandler.GetStatusHistory)
	reqpickup.Post("/otomatis", pickupHandler.CreateRequestPickup)
	reqpickup.Put("/:id/select-collector", pickupHandler.SelectCollector)
	reqpickup.Put("/pickup/:id/status", pickupHandler.UpdatePickupStatus)
	reqpickup.Put("/pickup/:id/item/update-actual", pickupHandler.UpdatePickupItemActualAmount)
}
