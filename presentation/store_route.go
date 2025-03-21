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

func StoreRouter(api fiber.Router) {

	storeRepo := repositories.NewStoreRepository(config.DB)
	storeService := services.NewStoreService(storeRepo)
	storeHandler := handler.NewStoreHandler(storeService)

	storeAPI := api.Group("/storerijig")

	storeAPI.Post("/create", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), storeHandler.CreateStore)
	storeAPI.Get("/getbyuser", middleware.AuthMiddleware, storeHandler.GetStoreByUserID)
	storeAPI.Put("/update/:store_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), storeHandler.UpdateStore)
	storeAPI.Delete("/delete/:store_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), storeHandler.DeleteStore)
}
