package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func APIKeyMiddleware(c *fiber.Ctx) error {

	apiKey := c.Get("x-api-key")

	expectedAPIKey := os.Getenv("API_KEY")

	if apiKey != expectedAPIKey {

		response := utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Invalid API Key",
			nil,
		)

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	return c.Next()
}
