package whatsapp

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func WhatsAppRouter(api fiber.Router) {

	whatsapp := api.Group("/whatsapp")

	whatsapp.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())

	whatsapp.Post("/generate-qr", GenerateQRHandler)
	whatsapp.Get("/status", CheckLoginStatusHandler)
	whatsapp.Post("/logout", WhatsAppLogoutHandler)

	messaging := whatsapp.Group("/message")
	messaging.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
	messaging.Post("/send", ValidateSendMessageRequest, SendMessageHandler)

	management := whatsapp.Group("/management")
	management.Use(middleware.AuthMiddleware(), middleware.RequireAdminRole())
	management.Get("/device-info", GetDeviceInfoHandler)
	management.Get("/health", HealthCheckHandler)

	api.Get("/whatsapp/ping", func(c *fiber.Ctx) error {
		return utils.Success(c, "WhatsApp service is running")
	})
}
