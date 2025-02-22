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
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	userProfile, err := h.UserProfileService.GetUserProfile(userID)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, userProfile, "User profile retrieved successfully")
}

func (h *UserProfileHandler) UpdateUserProfile(c *fiber.Ctx) error {
	var updateData dto.UpdateUserDTO
	if err := c.BodyParser(&updateData); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	errors, valid := updateData.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userResponse, err := h.UserProfileService.UpdateUserProfile(userID, updateData)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.SuccessResponse(c, userResponse, "User profile updated successfully")
}

func (h *UserProfileHandler) UpdateUserPassword(c *fiber.Ctx) error {
	var passwordData dto.UpdatePasswordDTO
	if err := c.BodyParser(&passwordData); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	errors, valid := passwordData.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	message, err := h.UserProfileService.UpdateUserPassword(userID, passwordData)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, message)
}
func (h *UserProfileHandler) UpdateUserAvatar(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "No avatar file uploaded")
	}

	message, err := h.UserProfileService.UpdateUserAvatar(userID, file)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, message)
}
