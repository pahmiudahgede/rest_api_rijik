package middleware

import (
	"fmt"
	"log"
	"os"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: No token provided")
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error parsing token: %v", err)
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil || claims["device_id"] == nil {
		log.Println("Invalid token claims")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token claims")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Unexpected signing method")
	}

	userID := claims["sub"].(string)
	deviceID := claims["device_id"].(string)

	sessionKey := fmt.Sprintf("session:%s:%s", userID, deviceID)
	sessionData, err := utils.GetJSONData(sessionKey)
	if err != nil || sessionData == nil {
		log.Printf("Session expired or invalid for userID: %s, deviceID: %s", userID, deviceID)
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Session expired or invalid")
	}

	roleID, roleOK := sessionData["roleID"].(string)
	roleName, roleNameOK := sessionData["roleName"].(string)
	if !roleOK || !roleNameOK {
		log.Println("Invalid session data for userID:", userID)
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid session data")
	}

	c.Locals("userID", userID)
	c.Locals("roleID", roleID)
	c.Locals("roleName", roleName)
	c.Locals("device_id", deviceID)

	return c.Next()
}
