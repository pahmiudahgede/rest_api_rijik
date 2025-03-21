package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request dto.RegisterRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"invalid request body"}})
	}

	if errors, valid := request.Validate(); !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	err := h.authService.RegisterUser(&request)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "OTP has been sent to your phone")
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var request dto.VerifyOTPRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"invalid request body"}})
	}

	err := h.authService.VerifyOTP(&request)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "User successfully registered")
}
