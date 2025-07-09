// internal/admin/approval_handler.go
package admin

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type ApprovalHandler struct {
	service   ApprovalService
}

func NewApprovalHandler(service ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{
		service:   service,
	}
}

// GetPendingUsers menampilkan daftar pengguna yang menunggu persetujuan
// @Summary Get pending users for approval
// @Description Retrieve list of users (pengelola/pengepul) waiting for admin approval with filtering and pagination
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param role query string false "Filter by role" Enums(pengelola, pengepul)
// @Param status query string false "Filter by status" Enums(awaiting_approval, pending) default(awaiting_approval)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Success 200 {object} utils.Response{data=PendingUsersListResponse} "List of pending users"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/pending [get]
func (h *ApprovalHandler) GetPendingUsers(c *fiber.Ctx) error {
	// Parse query parameters
	var req GetPendingUsersRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid query parameters: "+err.Error())
	}

	// Validate request
	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	// }

	// Call service
	result, err := h.service.GetPendingUsers(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to get pending users: "+err.Error())
	}

	return utils.SuccessWithData(c, "Pending users retrieved successfully", result)
}

// GetUserApprovalDetails menampilkan detail lengkap pengguna untuk approval
// @Summary Get user approval details
// @Description Get detailed information of a specific user for approval decision
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Success 200 {object} utils.Response{data=PendingUserResponse} "User approval details"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/{user_id}/approval-details [get]
func (h *ApprovalHandler) GetUserApprovalDetails(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required")
	}

	// Validate UUID format
	// if err := h.validator.Var(userID, "uuid"); err != nil {
	// 	return utils.BadRequest(c, "Invalid user ID format")
	// }

	// Call service
	result, err := h.service.GetUserApprovalDetails(c.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to get user details: "+err.Error())
	}

	return utils.SuccessWithData(c, "User approval details retrieved successfully", result)
}

// ProcessApprovalAction memproses aksi approval (approve/reject) untuk satu user
// @Summary Process approval action
// @Description Approve or reject a user registration
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param request body ApprovalActionRequest true "Approval action request"
// @Success 200 {object} utils.Response{data=ApprovalActionResponse} "Approval processed successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/approval-action [post]
func (h *ApprovalHandler) ProcessApprovalAction(c *fiber.Ctx) error {
	// Get admin ID from context
	adminClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Admin authentication required")
	}

	// Parse request body
	var req ApprovalActionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body: "+err.Error())
	}

	// Validate request
	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	// }

	// Call service
	result, err := h.service.ProcessApprovalAction(c.Context(), &req, adminClaims.UserID)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to process approval: "+err.Error())
	}

	actionMessage := "User approved successfully"
	if req.Action == "reject" {
		actionMessage = "User rejected successfully"
	}

	return utils.SuccessWithData(c, actionMessage, result)
}

// BulkProcessApproval memproses aksi approval untuk multiple users sekaligus
// @Summary Bulk process approval actions
// @Description Approve or reject multiple users at once
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param request body BulkApprovalRequest true "Bulk approval request"
// @Success 200 {object} utils.Response{data=BulkApprovalResponse} "Bulk approval processed"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/bulk-approval [post]
func (h *ApprovalHandler) BulkProcessApproval(c *fiber.Ctx) error {
	// Get admin ID from context
	adminClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Admin authentication required")
	}

	// Parse request body
	var req BulkApprovalRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body: "+err.Error())
	}

	// Validate request
	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	// }

	// Call service
	result, err := h.service.BulkProcessApproval(c.Context(), &req, adminClaims.UserID)
	if err != nil {
		return utils.InternalServerError(c, "Failed to process bulk approval: "+err.Error())
	}

	actionMessage := "Bulk approval processed successfully"
	if req.Action == "reject" {
		actionMessage = "Bulk rejection processed successfully"
	}

	return utils.SuccessWithData(c, actionMessage, result)
}

// ApproveUser endpoint khusus untuk approve satu user (shortcut)
// @Summary Approve user
// @Description Approve a user registration (shortcut endpoint)
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param notes body string false "Optional approval notes"
// @Success 200 {object} utils.Response{data=ApprovalActionResponse} "User approved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/{user_id}/approve [post]
func (h *ApprovalHandler) ApproveUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required")
	}

	// Validate UUID format
	// if err := h.validator.Var(userID, "uuid"); err != nil {
	// 	return utils.BadRequest(c, "Invalid user ID format")
	// }

	// Get admin ID from context
	adminClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Admin authentication required")
	}

	// Parse optional notes from body
	var body struct {
		Notes string `json:"notes"`
	}
	c.BodyParser(&body) // Ignore error as notes are optional

	// Call service
	result, err := h.service.ApproveUser(c.Context(), userID, adminClaims.UserID, body.Notes)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to approve user: "+err.Error())
	}

	return utils.SuccessWithData(c, "User approved successfully", result)
}

// RejectUser endpoint khusus untuk reject satu user (shortcut)
// @Summary Reject user
// @Description Reject a user registration (shortcut endpoint)
// @Tags Admin - User Approval
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param notes body string false "Optional rejection notes"
// @Success 200 {object} utils.Response{data=ApprovalActionResponse} "User rejected successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin role required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Security Bearer
// @Router /admin/users/{user_id}/reject [post]
func (h *ApprovalHandler) RejectUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required")
	}

	// Validate UUID format
	// if err := h.validator.Var(userID, "uuid"); err != nil {
	// 	return utils.BadRequest(c, "Invalid user ID format")
	// }

	// Get admin ID from context
	adminClaims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Admin authentication required")
	}

	// Parse optional notes from body
	var body struct {
		Notes string `json:"notes"`
	}
	c.BodyParser(&body) // Ignore error as notes are optional

	// Call service
	result, err := h.service.RejectUser(c.Context(), userID, adminClaims.UserID, body.Notes)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to reject user: "+err.Error())
	}

	return utils.SuccessWithData(c, "User rejected successfully", result)
}