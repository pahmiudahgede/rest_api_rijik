package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func CollectorRouter(api fiber.Router) {
	cartRepo :=       repositories.NewCartRepository()
	// trashRepo repositories.TrashRepository

	pickupRepo := repositories.NewRequestPickupRepository()
	trashRepo := repositories.NewTrashRepository(config.DB)
	historyRepo := repositories.NewPickupStatusHistoryRepository()
	cartService := services.NewCartService(cartRepo, trashRepo)

	pickupService := services.NewRequestPickupService(trashRepo, pickupRepo, cartService, historyRepo)
	pickupHandler := handler.NewRequestPickupHandler(pickupService)
	collectorRepo := repositories.NewCollectorRepository()

	collectorService := services.NewCollectorService(collectorRepo, trashRepo)
	collectorHandler := handler.NewCollectorHandler(collectorService)

	collectors := api.Group("/collectors")
	collectors.Use(middleware.AuthMiddleware)

	collectors.Post("/", collectorHandler.CreateCollector)
	collectors.Post("/:id/trash", collectorHandler.AddTrashToCollector)
	collectors.Get("/:id", collectorHandler.GetCollectorByID)
	collectors.Get("/", collectorHandler.GetCollectorByUserID)
	collectors.Get("/pickup/assigned-to-me", pickupHandler.GetAssignedPickup)

	collectors.Patch("/:id", collectorHandler.UpdateCollector)
	collectors.Patch("/:id/trash", collectorHandler.UpdateTrash)
	collectors.Patch("/:id/job-status", collectorHandler.UpdateJobStatus)
	collectors.Delete("/trash/:id", collectorHandler.DeleteTrash)
}
