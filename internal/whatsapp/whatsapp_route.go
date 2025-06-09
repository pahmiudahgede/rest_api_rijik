package whatsapp

import (
	"github.com/gofiber/fiber/v2"
)

func WhatsAppRouter(api fiber.Router) {
	api.Get("/whatsapp-status", WhatsAppStatusHandler)
	api.Get("/whatsapp/pw=admin1234", WhatsAppQRPageHandler)
	api.Post("/logout/whastapp", WhatsAppLogoutHandler)
}
