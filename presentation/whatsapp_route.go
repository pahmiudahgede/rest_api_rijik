package presentation

import (
	"rijig/internal/handler"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func WhatsAppRouter(api fiber.Router) {
	api.Post("/logout/whastapp", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), handler.WhatsAppHandler)
}
