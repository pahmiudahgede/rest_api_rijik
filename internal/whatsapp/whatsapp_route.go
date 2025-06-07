package whatsapp

import (
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func WhatsAppRouter(api fiber.Router) {
	api.Post("/logout/whastapp", middleware.AuthMiddleware(), WhatsAppHandler)
}
