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

	requestRepo := repositories.NewRequestPickupRepository(config.DB)
	repoTrash := repositories.NewTrashRepository(config.DB)
	repoAddress := repositories.NewAddressRepository(config.DB)

	requestPickupServices := services.NewRequestPickupService(requestRepo, repoAddress, repoTrash)

	requestPickupHandler := handler.NewRequestPickupHandler(requestPickupServices)

	requestPickupAPI := api.Group("/requestpickup")
	requestPickupAPI.Use(middleware.AuthMiddleware)

	requestPickupAPI.Post("/", requestPickupHandler.CreateRequestPickup)
	// requestPickupAPI.Get("/:id", requestPickupHandler.GetRequestPickupByID)
	// requestPickupAPI.Get("/", requestPickupHandler.GetAllRequestPickups)
	// requestPickupAPI.Put("/:id", requestPickupHandler.UpdateRequestPickup)
	// requestPickupAPI.Delete("/:id", requestPickupHandler.DeleteRequestPickup)
}
