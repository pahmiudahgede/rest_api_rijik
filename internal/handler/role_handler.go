package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type RoleHandler struct {
	RoleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: roleService}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {

	roleID, ok := c.Locals("roleID").(string)
	if !ok || roleID != utils.RoleAdministrator {
		return utils.GenericResponse(c, fiber.StatusForbidden, "Forbidden: You don't have permission to access this resource")
	}

	roles, err := h.RoleService.GetRoles()
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, roles, "Roles fetched successfully")
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	roleID := c.Params("role_id")

	roleIDFromSession, ok := c.Locals("roleID").(string)
	if !ok || roleIDFromSession != utils.RoleAdministrator {
		return utils.GenericResponse(c, fiber.StatusForbidden, "Forbidden: You don't have permission to access this resource")
	}

	role, err := h.RoleService.GetRoleByID(roleID)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "role id tidak ditemukan")
	}

	return utils.SuccessResponse(c, role, "Role fetched successfully")
}
