package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func CollectorRouter(api fiber.Router) {
	repo := repositories.NewCollectorRepository(config.DB)
	repoReq := repositories.NewRequestPickupRepository(config.DB)
	repoAddress := repositories.NewAddressRepository(config.DB)
	colectorService := services.NewCollectorService(repo, repoReq, repoAddress)
	collectorHandler := handler.NewCollectorHandler(colectorService)

	collector := api.Group("/collector")
	collector.Use(middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RolePengepul))

	collector.Put("confirmrequest/:id", collectorHandler.ConfirmRequestPickup)

}
