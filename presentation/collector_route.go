package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	// "rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func CollectorRouter(api fiber.Router) {
	// repo := repositories.NewCollectorRepository(config.DB)
	// repoReq := repositories.NewRequestPickupRepository(config.DB)
	// repoAddress := repositories.NewAddressRepository(config.DB)
	// repoUser := repositories.NewUserProfilRepository(config.DB)
	// colectorService := services.NewCollectorService(repo, repoReq, repoAddress, repoUser)
	// collectorHandler := handler.NewCollectorHandler(colectorService)

	// collector := api.Group("/collector")
	// collector.Use(middleware.AuthMiddleware)

	// collector.Put("confirmrequest/:id", collectorHandler.ConfirmRequestPickup)
	// collector.Put("confirm-manual/request/:request_id", collectorHandler.ConfirmRequestManualPickup)
	// collector.Get("/avaible", collectorHandler.GetAvaibleCollector)

	// Middleware Auth dan Role

	// Inisialisasi repository dan service
	pickupRepo := repositories.NewRequestPickupRepository()
	trashRepo := repositories.NewTrashRepository(config.DB)
	historyRepo := repositories.NewPickupStatusHistoryRepository()
	cartService := services.NewCartService()

	pickupService := services.NewRequestPickupService(trashRepo, pickupRepo, cartService, historyRepo)
	pickupHandler := handler.NewRequestPickupHandler(pickupService)
	collectorRepo := repositories.NewCollectorRepository()
	// trashRepo := repositories.NewTrashRepository(config.DB)
	collectorService := services.NewCollectorService(collectorRepo, trashRepo)
	collectorHandler := handler.NewCollectorHandler(collectorService)

	// Group collector
	collectors := api.Group("/collectors")
	collectors.Use(middleware.AuthMiddleware)

	// === Collector routes ===
	collectors.Post("/", collectorHandler.CreateCollector)              // Create collector
	collectors.Post("/:id/trash", collectorHandler.AddTrashToCollector) // Add trash to collector
	collectors.Get("/:id", collectorHandler.GetCollectorByID)           // Get collector by ID
	collectors.Get("/", collectorHandler.GetCollectorByUserID) 
	collectors.Get("/pickup/assigned-to-me", pickupHandler.GetAssignedPickup) 
         // Get collector by ID
	collectors.Patch("/:id", collectorHandler.UpdateCollector)          // Update collector fields
	collectors.Patch("/:id/trash", collectorHandler.UpdateTrash)
	collectors.Patch("/:id/job-status", collectorHandler.UpdateJobStatus)
	collectors.Delete("/trash/:id", collectorHandler.DeleteTrash)
}
