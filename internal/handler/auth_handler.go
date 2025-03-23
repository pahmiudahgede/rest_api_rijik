package handler

import (
	"log"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) RegisterOrLoginHandler(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	if req.Phone == "" || req.RoleID == "" {
		return utils.ErrorResponse(c, "Phone number and role ID are required")
	}

	if err := h.authService.RegisterOrLogin(&req); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "OTP sent successfully")
}

func (h *AuthHandler) VerifyOTPHandler(c *fiber.Ctx) error {
	var req dto.VerifyOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	if req.OTP == "" {
		return utils.ErrorResponse(c, "OTP is required")
	}

	response, err := h.authService.VerifyOTP(&req)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, response, "Registration/Login successful")
}

func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "User is not logged in or invalid session")
	}

	phoneKey := "user_phone:" + userID
	phone, err := utils.GetStringData(phoneKey)
	if err != nil || phone == "" {

		log.Printf("Error retrieving phone from Redis for user %s: %v", userID, err)
		return utils.ErrorResponse(c, "Phone number is missing or invalid session data")
	}

	err = h.authService.Logout(userID, phone)
	if err != nil {

		log.Printf("Error during logout process for user %s: %v", userID, err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Logged out successfully")
}
