package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, exists := c.Locals("roleID").(string)
		if !exists || roleID == "" {
			return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Role not found in session")
		}

		for _, role := range allowedRoles {
			if role == roleID {
				return c.Next()
			}
		}

		return utils.GenericErrorResponse(c, fiber.StatusForbidden, "Access Denied: You don't have permission to access this resource")
	}
}
