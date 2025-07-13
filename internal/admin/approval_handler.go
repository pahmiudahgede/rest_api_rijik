package admin

import (
	"log"
	"rijig/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService AdminService
	validator    *validator.Validate
}

func NewAdminHandler(adminService AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		validator:    validator.New(),
	}
}

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	req, err := h.parseGetAllUsersRequest(c)
	if err != nil {
		log.Printf("Error parsing request parameters: %v", err)
		return utils.BadRequest(c, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("Validation error: %v", err)
		return utils.BadRequest(c, "Invalid request parameters")
	}

	validPage, validLimit, err := h.adminService.ValidatePaginationParams(req.Page, req.Limit)
	if err != nil {
		log.Printf("Pagination validation error: %v", err)
		return utils.BadRequest(c, err.Error())
	}

	req.Page = validPage
	req.Limit = validLimit

	data, page, limit, total, err := h.adminService.GetAllUsers(ctx, *req)
	if err != nil {
		log.Printf("Service error in GetAllUsers: %v", err)
		return utils.InternalServerError(c, "Failed to retrieve users")
	}

	hasPagination := page != nil && limit != nil
	message := h.adminService.GetMessage(req.Role, hasPagination, total)

	if hasPagination {
		return utils.SuccessWithPaginationAndTotal(c, message, data, *page, *limit, int(total))
	}

	return utils.SuccessWithTotal(c, message, data, int(total))
}

func (h *AdminHandler) UpdateRegistrationStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Params("userid")
	if userID == "" {
		log.Printf("Missing userid parameter")
		return utils.BadRequest(c, "User ID is required")
	}

	var req UpdateRegistrationStatusRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("Validation error for registration status update: %v", err)
		return utils.BadRequest(c, "Invalid action. Must be 'approved' or 'rejected'")
	}

	if err := h.adminService.UpdateRegistrationStatus(ctx, userID, req); err != nil {
		log.Printf("Service error in UpdateRegistrationStatus: %v", err)

		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}

		if err.Error() == "failed to update registration status" {
			return utils.InternalServerError(c, "Failed to update registration status")
		}

		return utils.BadRequest(c, err.Error())
	}

	message := h.generateUpdateMessage(req.Action)
	return utils.Success(c, message)
}

func (h *AdminHandler) GetUserStatistics(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.adminService.GetUserStatistics(ctx)
	if err != nil {
		log.Printf("Error getting user statistics: %v", err)
		return utils.InternalServerError(c, "Failed to get user statistics")
	}

	return utils.SuccessWithData(c, "Successfully retrieved user statistics", stats)
}

func (h *AdminHandler) parseGetAllUsersRequest(c *fiber.Ctx) (*GetAllUsersRequest, error) {
	req := &GetAllUsersRequest{}

	req.Role = c.Query("role")
	if req.Role == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "role parameter is required")
	}

	req.StatusReg = c.Query("statusreg")

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid page parameter")
		}
		req.Page = &page
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid limit parameter")
		}
		req.Limit = &limit
	}

	return req, nil
}

func (h *AdminHandler) generateUpdateMessage(action string) string {
	switch action {
	case "approved":
		return "User registration has been approved successfully"
	case "rejected":
		return "User registration has been rejected successfully"
	default:
		return "Registration status updated successfully"
	}
}

func (h *AdminHandler) ValidateAdminPermissions(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")
	if userRole != "admin" && userRole != "super_admin" {
		return utils.Forbidden(c, "Admin permissions required")
	}

	return nil
}

func (h *AdminHandler) GetAllUsersExport(c *fiber.Ctx) error {
	ctx := c.Context()

	req, err := h.parseGetAllUsersRequest(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	req.Page = nil
	req.Limit = nil

	data, _, _, total, err := h.adminService.GetAllUsers(ctx, *req)
	if err != nil {
		log.Printf("Error exporting users: %v", err)
		return utils.InternalServerError(c, "Failed to export users")
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=users_export.csv")

	return utils.SuccessWithTotal(c, "Users data for export", data, int(total))
}

func (h *AdminHandler) GetUsersByRole(c *fiber.Ctx) error {
	role := c.Params("role")
	if role == "" {
		return utils.BadRequest(c, "Role parameter is required")
	}

	req := GetAllUsersRequest{
		Role: role,
	}

	return h.handleGetUsersRequest(c, req)
}

func (h *AdminHandler) handleGetUsersRequest(c *fiber.Ctx, req GetAllUsersRequest) error {
	ctx := c.Context()

	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, "Invalid request parameters")
	}

	data, page, limit, total, err := h.adminService.GetAllUsers(ctx, req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve users")
	}

	hasPagination := page != nil && limit != nil
	message := h.adminService.GetMessage(req.Role, hasPagination, total)

	if hasPagination {
		return utils.SuccessWithPaginationAndTotal(c, message, data, *page, *limit, int(total))
	}

	return utils.SuccessWithTotal(c, message, data, int(total))
}

func (h *AdminHandler) HealthCheck(c *fiber.Ctx) error {
	return utils.SuccessWithData(c, "Admin module is healthy", fiber.Map{
		"module":  "admin",
		"status":  "healthy",
		"version": "1.0.0",
	})
}
