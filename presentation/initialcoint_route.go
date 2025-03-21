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

func InitialCointRoute(api fiber.Router) {

	initialCointRepo := repositories.NewInitialCointRepository(config.DB)
	initialCointService := services.NewInitialCointService(initialCointRepo)
	initialCointHandler := handler.NewInitialCointHandler(initialCointService)

	initialCoint := api.Group("/initialcoint")
	initialCoint.Use(middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator))

	initialCoint.Post("/create", initialCointHandler.CreateInitialCoint)
	initialCoint.Get("/getall", initialCointHandler.GetAllInitialCoints)
	initialCoint.Get("/get/:coin_id", initialCointHandler.GetInitialCointByID)
	initialCoint.Put("/update/:coin_id", initialCointHandler.UpdateInitialCoint)
	initialCoint.Delete("/delete/:coin_id", initialCointHandler.DeleteInitialCoint)
}
