package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("roleID").(string)
		if !ok || roleID == "" {
			log.Println("Unauthorized access: Role not found in session")
			return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Role not found in session")
		}

		for _, role := range allowedRoles {
			if role == roleID {
				return c.Next()
			}
		}

		log.Println("Access denied for role:", roleID)
		return utils.GenericErrorResponse(c, fiber.StatusForbidden, "Access Denied: You don't have permission to access this resource")
	}
}
