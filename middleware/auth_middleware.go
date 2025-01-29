package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: No token provided")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token claims")
	}

	sessionKey := "session:" + claims["sub"].(string)
	sessionData, err := utils.GetJSONData(sessionKey)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Session expired or invalid")
	}

	c.Locals("userID", sessionData["userID"])
	c.Locals("roleID", sessionData["roleID"])
	c.Locals("roleName", sessionData["roleName"])

	return c.Next()
}
