package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type UserProfileHandler struct {
	UserProfileService services.UserProfileService
}

func NewUserProfileHandler(userProfileService services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{UserProfileService: userProfileService}
}

func (h *UserProfileHandler) GetUserProfile(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	userProfile, err := h.UserProfileService.GetUserProfile(userID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.LogResponse(c, userProfile, "User profile retrieved successfully")
}

func (h *UserProfileHandler) UpdateUserProfile(c *fiber.Ctx) error {
	var updateData dto.UpdateUserDTO
	if err := c.BodyParser(&updateData); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	errors, valid := updateData.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userResponse, err := h.UserProfileService.UpdateUserProfile(userID, updateData)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.LogResponse(c, userResponse, "User profile updated successfully")
}

func (h *UserProfileHandler) UpdateUserPassword(c *fiber.Ctx) error {
	var passwordData dto.UpdatePasswordDTO
	if err := c.BodyParser(&passwordData); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	errors, valid := passwordData.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userResponse, err := h.UserProfileService.UpdateUserPassword(userID, passwordData)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.LogResponse(c, userResponse, "Password updated successfully")
}
