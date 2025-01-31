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

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token claims")
	}

	userID := claims["sub"].(string)
	if userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid user session")
	}

	sessionKey := "session:" + userID
	sessionData, err := utils.GetJSONData(sessionKey)
	if err != nil || sessionData == nil {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Session expired or invalid")
	}

	roleID, roleOK := sessionData["roleID"].(string)
	roleName, roleNameOK := sessionData["roleName"].(string)
	if !roleOK || !roleNameOK {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid session data")
	}

	c.Locals("userID", userID)
	c.Locals("roleID", roleID)
	c.Locals("roleName", roleName)

	return c.Next()
}
