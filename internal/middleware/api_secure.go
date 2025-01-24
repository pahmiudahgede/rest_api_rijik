package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func APIKeyMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("x-api-key")
	expectedAPIKey := os.Getenv("API_KEY")

	if apiKey != expectedAPIKey {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Invalid API Key",
			nil,
		))
	}

	return c.Next()
}

func RateLimitMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("x-api-key")
	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"API Key is missing",
			nil,
		))
	}

	ctx := context.Background()
	rateLimitKey := fmt.Sprintf("rate_limit:%s", apiKey)

	count, _ := config.RedisClient.Incr(ctx, rateLimitKey).Result()
	if count > 100 {
		return c.Status(fiber.StatusTooManyRequests).JSON(utils.FormatResponse(
			fiber.StatusTooManyRequests,
			"Rate limit exceeded",
			nil,
		))
	}

	config.RedisClient.Expire(ctx, rateLimitKey, time.Minute)

	return c.Next()
}
