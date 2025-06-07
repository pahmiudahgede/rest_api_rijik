package presentation

import (
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func PickupMatchingRouter(api fiber.Router) {
	pickupRepo := repositories.NewRequestPickupRepository()
	collectorRepo := repositories.NewCollectorRepository()
	service := services.NewPickupMatchingService(pickupRepo, collectorRepo)
	handler := handler.NewPickupMatchingHandler(service)

	manual := api.Group("/pickup/manual")
	manual.Use(middleware.AuthMiddleware())
	manual.Get("/:pickupID/nearby-collectors", handler.GetNearbyCollectorsForPickup)

	auto := api.Group("/pickup/otomatis")
	auto.Use(middleware.AuthMiddleware())
	auto.Get("/available-requests", handler.GetAvailablePickupForCollector)
}
