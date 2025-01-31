package handler

import (
	"github.com/gofiber/fiber/v2"
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
