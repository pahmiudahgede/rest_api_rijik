package role

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService RoleService
}

func NewRoleHandler(roleService RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {

	if _, err := middleware.GetUserFromContext(c); err != nil {
		return utils.Unauthorized(c, "Unauthorized access")
	}

	roles, err := h.roleService.GetRoles(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch roles")
	}

	return utils.SuccessWithData(c, "Roles fetched successfully", roles)
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {

	if _, err := middleware.GetUserFromContext(c); err != nil {
		return utils.Unauthorized(c, "Unauthorized access")
	}

	roleID := c.Params("role_id")
	if roleID == "" {
		return utils.BadRequest(c, "Role ID is required")
	}

	role, err := h.roleService.GetRoleByID(c.Context(), roleID)
	if err != nil {

		return utils.NotFound(c, "Role not found")
	}

	return utils.SuccessWithData(c, "Role fetched successfully", role)
}
