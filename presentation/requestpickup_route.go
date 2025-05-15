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
	// repo        repositories.RequestPickupRepository
	// repoColl    repositories.CollectorRepository
	// repoAddress repositories.AddressRepository
	// repoTrash   repositories.TrashRepository
	// repoUser    repositories.UserProfilRepository

	requestRepo := repositories.NewRequestPickupRepository(config.DB)
	repoColl := repositories.NewCollectorRepository(config.DB)
	repoAddress := repositories.NewAddressRepository(config.DB)
	Trashrepo := repositories.NewTrashRepository(config.DB)
	repouser := repositories.NewUserProfilRepository(config.DB)
	// collectorRepo := repositories.NewCollectorRepository(config.DB)

	requestPickupServices := services.NewRequestPickupService(requestRepo, repoColl, repoAddress, Trashrepo, repouser)
	// collectorService := services.NewCollectorService(collectorRepo, requestRepo, repoAddress)
	// service services.RequestPickupService,
	// collectorService services.CollectorService

	requestPickupHandler := handler.NewRequestPickupHandler(requestPickupServices)

	requestPickupAPI := api.Group("/requestpickup")
	requestPickupAPI.Use(middleware.AuthMiddleware)

	requestPickupAPI.Post("/", requestPickupHandler.CreateRequestPickup)
	// requestPickupAPI.Get("/get", middleware.AuthMiddleware, requestPickupHandler.GetAutomaticRequestByUser)
	requestPickupAPI.Get("/get-allrequest", requestPickupHandler.GetRequestPickups)
	requestPickupAPI.Patch("/select-collector", requestPickupHandler.AssignCollectorToRequest)
	// requestPickupAPI.Get("/:id", requestPickupHandler.GetRequestPickupByID)
	// requestPickupAPI.Get("/", requestPickupHandler.GetAllRequestPickups)
	// requestPickupAPI.Put("/:id", requestPickupHandler.UpdateRequestPickup)
	// requestPickupAPI.Delete("/:id", requestPickupHandler.DeleteRequestPickup)
}
