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
