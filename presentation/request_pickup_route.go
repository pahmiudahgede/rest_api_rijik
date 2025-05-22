package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func RequestPickupRouter(api fiber.Router) {
	pickupRepo := repositories.NewRequestPickupRepository()
	historyRepo := repositories.NewPickupStatusHistoryRepository()
	trashRepo := repositories.NewTrashRepository(config.DB)
	cartService := services.NewCartService()
	historyService := services.NewPickupStatusHistoryService(historyRepo)

	pickupService := services.NewRequestPickupService(trashRepo, pickupRepo, cartService, historyRepo)
	pickupHandler := handler.NewRequestPickupHandler(pickupService)
	statuspickupHandler := handler.NewPickupStatusHistoryHandler(historyService)

	reqpickup := api.Group("/reqpickup")
	reqpickup.Use(middleware.AuthMiddleware)

	reqpickup.Post("/manual", pickupHandler.CreateRequestPickup)
	reqpickup.Get("/pickup/:id/history", statuspickupHandler.GetStatusHistory)
	reqpickup.Post("/otomatis", pickupHandler.CreateRequestPickup)
	reqpickup.Put("/:id/select-collector", pickupHandler.SelectCollector)
	reqpickup.Put("/pickup/:id/status", pickupHandler.UpdatePickupStatus)
	reqpickup.Put("/pickup/:id/item/update-actual", pickupHandler.UpdatePickupItemActualAmount)
}
