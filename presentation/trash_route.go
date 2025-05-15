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

func TrashRouter(api fiber.Router) {
	trashRepo := repositories.NewTrashRepository(config.DB)
	trashService := services.NewTrashService(trashRepo)
	trashHandler := handler.NewTrashHandler(trashService)

	trashAPI := api.Group("/trash")
	trashAPI.Use(middleware.AuthMiddleware)

	trashAPI.Post("/category", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.CreateCategory)
	trashAPI.Post("/category/detail", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.AddDetailToCategory)
	trashAPI.Get("/categories", trashHandler.GetCategories)
	trashAPI.Get("/category/:category_id", trashHandler.GetCategoryByID)
	trashAPI.Get("/detail/:detail_id", trashHandler.GetTrashDetailByID)

	trashAPI.Patch("/category/:category_id", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.UpdateCategory)
	trashAPI.Put("/detail/:detail_id", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.UpdateDetail)

	trashAPI.Delete("/category/:category_id", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.DeleteCategory)
	trashAPI.Delete("/detail/:detail_id", middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.DeleteDetail)
}
