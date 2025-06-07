package middleware

import (
	"os"

	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("x-api-key")
	if apiKey == "" {
		return utils.Unauthorized(c, "Unauthorized: API key is required")
	}

	validAPIKey := os.Getenv("API_KEY")
	if apiKey != validAPIKey {
		return utils.Unauthorized(c, "Unauthorized: Invalid API key")
	}

	return c.Next()
}
