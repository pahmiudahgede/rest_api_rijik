package whatsapp

import (
	"log"
	"rijig/config"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func WhatsAppHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "User is not logged in or invalid session")
	}

	err := config.LogoutWhatsApp()
	if err != nil {
		log.Printf("Error during logout process for user %s: %v", userID, err)
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Logged out successfully")
}
