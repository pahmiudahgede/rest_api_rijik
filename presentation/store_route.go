package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func StoreRouter(api fiber.Router) {

	storeRepo := repositories.NewStoreRepository(config.DB)
	storeService := services.NewStoreService(storeRepo)
	storeHandler := handler.NewStoreHandler(storeService)

	storeAPI := api.Group("/storerijig")
	storeAPI.Post("/create", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), storeHandler.CreateStore)
}
