package handler

import (
	"log"
	"rijig/config"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func WhatsAppHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "User is not logged in or invalid session")
	}

	err := config.LogoutWhatsApp()
	if err != nil {
		log.Printf("Error during logout process for user %s: %v", userID, err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Logged out successfully")
}
