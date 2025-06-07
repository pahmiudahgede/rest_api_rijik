package presentation

import (
	"rijig/config"
	"rijig/internal/about"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func AboutRouter(api fiber.Router) {
	aboutRepo := repositories.NewAboutRepository(config.DB)
	aboutService := services.NewAboutService(aboutRepo)
	aboutHandler := about.NewAboutHandler(aboutService)

	aboutRoutes := api.Group("/about")
	aboutRoutes.Use(middleware.AuthMiddleware())

	aboutRoutes.Get("/", aboutHandler.GetAllAbout)
	aboutRoutes.Get("/:id", aboutHandler.GetAboutByID)
	aboutRoutes.Post("/", aboutHandler.CreateAbout) // admin
	aboutRoutes.Put("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.UpdateAbout)
	aboutRoutes.Delete("/:id", aboutHandler.DeleteAbout) // admin

	aboutDetailRoutes := api.Group("/about-detail")
	aboutDetailRoutes.Use(middleware.AuthMiddleware())
	aboutDetailRoute := api.Group("/about-detail")
	aboutDetailRoute.Get("/:id", aboutHandler.GetAboutDetailById)
	aboutDetailRoutes.Post("/", aboutHandler.CreateAboutDetail) // admin
	aboutDetailRoutes.Put("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.UpdateAboutDetail)
	aboutDetailRoutes.Delete("/:id", middleware.RequireRoles(utils.RoleAdministrator), aboutHandler.DeleteAboutDetail)
}
