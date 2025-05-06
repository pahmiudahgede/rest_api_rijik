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

func AboutRouter(api fiber.Router) {

	aboutRepo := repositories.NewAboutRepository(config.DB)
	aboutService := services.NewAboutService(aboutRepo)
	aboutHandler := handler.NewAboutHandler(aboutService)

	aboutRoutes := api.Group("/about")
	aboutRoute := api.Group("/about")
	aboutRoutes.Use(middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator))

	aboutRoute.Get("/", aboutHandler.GetAllAbout)
	aboutRoute.Get("/:id", aboutHandler.GetAboutByID)
	aboutRoutes.Post("/", aboutHandler.CreateAbout)
	aboutRoutes.Put("/:id", aboutHandler.UpdateAbout)
	aboutRoutes.Delete("/:id", aboutHandler.DeleteAbout)

	aboutDetailRoutes := api.Group("/about-detail")
	aboutDetailRoutes.Use(middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator))
	aboutDetailRoute := api.Group("/about-detail")
	aboutDetailRoute.Get("/:id", aboutHandler.GetAboutDetailById)
	aboutDetailRoutes.Post("/", aboutHandler.CreateAboutDetail)
	aboutDetailRoutes.Put("/:id", aboutHandler.UpdateAboutDetail)
	aboutDetailRoutes.Delete("/:id", aboutHandler.DeleteAboutDetail)
}
