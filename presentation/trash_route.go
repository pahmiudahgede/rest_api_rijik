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

func TrashRouter(api fiber.Router) {

	trashRepo := repositories.NewTrashRepository(config.DB)
	trashService := services.NewTrashService(trashRepo)
	trashHandler := handler.NewTrashHandler(trashService)

	trashAPI := api.Group("/trash")

	trashAPI.Post("/createcategory-with-details", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.CreateCategory)
	trashAPI.Post("/add-detailinexistingcategory", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.AddDetailToCategory)
	trashAPI.Get("/allcategory", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), trashHandler.GetCategories)
	trashAPI.Get("/category/:category_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), trashHandler.GetCategoryByID)
	trashAPI.Get("/detail/:detail_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), trashHandler.GetTrashDetailByID)

	trashAPI.Patch("/update-category/:category_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.UpdateCategory)
	trashAPI.Put("/update-detail/:detail_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.UpdateDetail)

	trashAPI.Delete("delete-category/:category_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.DeleteCategory)
	trashAPI.Delete("delete-detail/:detail_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), trashHandler.DeleteDetail)
}
