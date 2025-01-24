package middleware

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func RoleRequired(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Token is required",
				nil,
			))
		}

		ctx := context.Background()
		cachedToken, err := config.RedisClient.Get(ctx, "auth_token:"+tokenString).Result()
		if err != nil || cachedToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			))
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("API_KEY")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Invalid token claims",
				nil,
			))
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Missing or invalid user ID in token",
				nil,
			))
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"Missing or invalid role in token",
				nil,
			))
		}

		c.Locals("userID", userID)
		c.Locals("role", role)

		if !containsRole(roles, role) {
			return c.Status(fiber.StatusForbidden).JSON(utils.FormatResponse(
				fiber.StatusForbidden,
				"You do not have permission to access this resource",
				nil,
			))
		}

		return c.Next()
	}
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Missing or invalid token",
			nil,
		))
	}

	ctx := context.Background()
	cachedToken, err := config.RedisClient.Get(ctx, "auth_token:"+tokenString).Result()
	if err != nil || cachedToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Invalid or expired token",
			nil,
		))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_KEY")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Invalid or expired token",
			nil,
		))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Invalid token claims",
			nil,
		))
	}

	userID := claims["sub"].(string)
	c.Locals("userID", userID)

	config.RedisClient.Expire(ctx, "auth_token:"+tokenString, time.Hour*24)

	return c.Next()
}