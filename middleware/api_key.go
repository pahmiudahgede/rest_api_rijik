package middleware

import (
	"os"

	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("x-api-key")
	if apiKey == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: API key is required")
	}

	validAPIKey := os.Getenv("API_KEY")
	if apiKey != validAPIKey {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid API key")
	}

	return c.Next()
}
