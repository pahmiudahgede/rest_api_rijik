package requestpickup

import (
	"rijig/config"
	"rijig/internal/collector"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func PickupMatchingRouter(api fiber.Router) {
	pickupRepo := NewRequestPickupRepository()
	collectorRepo := collector.NewCollectorRepository(config.DB)
	service := NewPickupMatchingService(pickupRepo, collectorRepo)
	handler := NewPickupMatchingHandler(service)

	manual := api.Group("/pickup/manual")
	manual.Use(middleware.AuthMiddleware)
	manual.Get("/:pickupID/nearby-collectors", handler.GetNearbyCollectorsForPickup)

	auto := api.Group("/pickup/otomatis")
	auto.Use(middleware.AuthMiddleware)
	auto.Get("/available-requests", handler.GetAvailablePickupForCollector)
}
