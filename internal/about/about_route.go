package about

import (
	"rijig/config"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func AboutRouter(api fiber.Router) {
	aboutRepo := NewAboutRepository(config.DB)
	aboutService := NewAboutService(aboutRepo)
	aboutHandler := NewAboutHandler(aboutService)

	aboutRoutes := api.Group("/about")
	aboutRoutes.Use(middleware.AuthMiddleware())

	aboutRoutes.Get("/", aboutHandler.GetAllAbout)
	aboutRoutes.Get("/:id", aboutHandler.GetAboutByID)
	aboutRoutes.Post("/", aboutHandler.CreateAbout)
	aboutRoutes.Put("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.UpdateAbout)
	aboutRoutes.Delete("/:id", aboutHandler.DeleteAbout)

	aboutDetailRoutes := api.Group("/about-detail")
	aboutDetailRoutes.Use(middleware.AuthMiddleware())
	aboutDetailRoute := api.Group("/about-detail")
	aboutDetailRoute.Get("/:id", aboutHandler.GetAboutDetailById)
	aboutDetailRoutes.Post("/", aboutHandler.CreateAboutDetail)
	aboutDetailRoutes.Put("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.UpdateAboutDetail)
	aboutDetailRoutes.Delete("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.DeleteAboutDetail)
}
