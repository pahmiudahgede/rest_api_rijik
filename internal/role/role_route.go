package role

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoleRouter(api fiber.Router) {
	roleRepo := NewRoleRepository(config.DB)
	roleService := NewRoleService(roleRepo)
	roleHandler := NewRoleHandler(roleService)
	roleRoute := api.Group("/role", middleware.AuthMiddleware())

	roleRoute.Get("/", roleHandler.GetRoles)
	roleRoute.Get("/:role_id", roleHandler.GetRoleByID)
}
