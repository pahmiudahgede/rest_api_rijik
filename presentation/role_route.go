// presentation/role_route.go
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

func RoleRouter(api fiber.Router) {
	roleRepo := repositories.NewRoleRepository(config.DB)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	api.Get("/roles", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), roleHandler.GetRoles)
	api.Get("/role/:role_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), roleHandler.GetRoleByID)
}
