package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
)

func RoleRouter(api fiber.Router) {
	roleRepo := repositories.NewRoleRepository(config.DB)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	api.Get("/roles", roleHandler.GetRoles)
	api.Get("/role/:role_id", roleHandler.GetRoleByID)
}
