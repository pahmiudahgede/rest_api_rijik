package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func APIKeyMiddleware(c *fiber.Ctx) error {

	apiKey := c.Get("x-api-key")

	validAPIKey := os.Getenv("API_KEY")

	if apiKey != validAPIKey {
		log.Printf("Invalid API Key: %s", apiKey)

		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: api key yang anda masukkan tidak valid")
	}

	return c.Next()
}
