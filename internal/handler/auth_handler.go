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
	return &AuthHandler{authService}
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	if errors, valid := req.Validate(); !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	err := h.authService.RegisterUser(&req)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Kode OTP telah dikirimkan ke nomor WhatsApp anda")
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req dto.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	response, err := h.authService.VerifyOTP(&req)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, response, "Registration successful")
}
