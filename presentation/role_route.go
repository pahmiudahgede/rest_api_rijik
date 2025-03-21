package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"

	"github.com/gofiber/fiber/v2"
)

func RoleRouter(api fiber.Router) {
	roleRepo := repositories.NewRoleRepository(config.DB)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	api.Get("/roles", roleHandler.GetRoles)
	api.Get("/role/:role_id", roleHandler.GetRoleByID)
}
