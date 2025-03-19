package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request dto.RegisterRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	if errors, valid := request.Validate(); !valid {
		return c.Status(400).JSON(errors)
	}

	_, err := h.AuthService.RegisterUser(request)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"meta": fiber.Map{
			"status":  201,
			"message": "The input register from the user has been successfully recorded. Please check the otp code sent to your number.",
		},
	})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var request struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	err := h.AuthService.VerifyOTP(request.Phone, request.OTP)
	if err != nil {
		return c.Status(400).JSON(dto.Response{
			Meta: dto.MetaResponse{
				Status:  400,
				Message: "Invalid OTP",
			},
			Data: nil,
		})
	}

	user, err := h.AuthService.GetUserByPhone(request.Phone)
	if err != nil {
		return c.Status(500).SendString("Error retrieving user")
	}
	if user == nil {
		return c.Status(404).SendString("User not found")
	}

	token, err := h.AuthService.GenerateJWT(user)
	if err != nil {
		return c.Status(500).SendString("Error generating token")
	}

	response := dto.Response{
		Meta: dto.MetaResponse{
			Status:  200,
			Message: "OTP yang dimasukkan valid",
		},
		Data: &dto.UserDataResponse{
			UserID:   user.ID,
			UserRole: user.Role.RoleName,
			Token:    token,
		},
	}

	return c.Status(200).JSON(response)
}
