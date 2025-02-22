package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if len(allowedRoles) == 0 {
			return utils.GenericResponse(c, fiber.StatusForbidden, "Forbidden: No roles specified")
		}

		roleID, ok := c.Locals("roleID").(string)
		if !ok || roleID == "" {
			return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: Role not found")
		}

		for _, role := range allowedRoles {
			if role == roleID {
				return c.Next()
			}
		}

		return utils.GenericResponse(c, fiber.StatusForbidden, "Access Denied: You don't have permission to access this resource")
	}
}
