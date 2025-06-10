package userprofile

import (
	"context"
	"log"
	"rijig/middleware"
	"rijig/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserProfileHandler struct {
	service UserProfileService
}

func NewUserProfileHandler(service UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		service: service,
	}
}

func (h *UserProfileHandler) GetUserProfile(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userProfile, err := h.service.GetUserProfile(ctx, claims.UserID)
	if err != nil {
		if strings.Contains(err.Error(), ErrUserNotFound.Error()) {
			return utils.NotFound(c, "User profile not found")
		}

		log.Printf("Error getting user profile: %v", err)
		return utils.InternalServerError(c, "Failed to retrieve user profile")
	}

	return utils.SuccessWithData(c, "User profile retrieved successfully", userProfile)
}

func (h *UserProfileHandler) UpdateUserProfile(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req RequestUserProfileDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if validationErrors, isValid := req.ValidateRequestUserProfileDTO(); !isValid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	updatedProfile, err := h.service.UpdateRegistUserProfile(ctx, claims.UserID, claims.DeviceID, &req)
	if err != nil {

		if strings.Contains(err.Error(), "user not found") {
			return utils.NotFound(c, "User not found")
		}

		log.Printf("Error updating user profile: %v", err)
		return utils.InternalServerError(c, "Failed to update user profile")
	}

	return utils.SuccessWithData(c, "User profile updated successfully", updatedProfile)
}
