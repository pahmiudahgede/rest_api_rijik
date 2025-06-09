package authentication

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthenticationHandler struct {
	service AuthenticationService
}

func NewAuthenticationHandler(service AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{service}
}

func (h *AuthenticationHandler) RefreshToken(c *fiber.Ctx) error {
	deviceID := c.Get("X-Device-ID")
	if deviceID == "" {
		return utils.BadRequest(c, "Device ID is required")
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if body.RefreshToken == "" {
		return utils.BadRequest(c, "Refresh token is required")
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized or missing user ID")
	}

	tokenData, err := utils.RefreshAccessToken(userID, deviceID, body.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.SuccessWithData(c, "Token refreshed successfully", tokenData)

}

func (h *AuthenticationHandler) GetMe(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)
	deviceID, _ := c.Locals("device_id").(string)
	regStatus, _ := c.Locals("registration_status").(string)

	data := fiber.Map{
		"user_id":             userID,
		"role":                role,
		"device_id":           deviceID,
		"registration_status": regStatus,
	}

	return utils.SuccessWithData(c, "User session data retrieved", data)

}

func (h *AuthenticationHandler) Login(c *fiber.Ctx) error {

	var req LoginAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if errs, ok := req.ValidateLoginAdminRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta": fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Validation failed",
			},
			"errors": errs,
		})
	}

	res, err := h.service.LoginAdmin(c.Context(), &req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.SuccessWithData(c, "Login successful", res)

}

func (h *AuthenticationHandler) Register(c *fiber.Ctx) error {

	var req RegisterAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if errs, ok := req.ValidateRegisterAdminRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta": fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "periksa lagi inputan",
			},
			"errors": errs,
		})
	}

	err := h.service.RegisterAdmin(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Registration successful, Please login")
}

func (h *AuthenticationHandler) RequestOtpHandler(c *fiber.Ctx) error {
	var req LoginorRegistRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if errs, ok := req.ValidateLoginorRegistRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta": fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Input tidak valid",
			},
			"errors": errs,
		})
	}

	_, err := h.service.SendLoginOTP(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "OTP sent successfully")
}

func (h *AuthenticationHandler) VerifyOtpHandler(c *fiber.Ctx) error {
	var req VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if errs, ok := req.ValidateVerifyOtpRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta":   fiber.Map{"status": fiber.StatusBadRequest, "message": "Validation error"},
			"errors": errs,
		})
	}

	stepResp, err := h.service.VerifyLoginOTP(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "OTP verified successfully", stepResp)
}

func (h *AuthenticationHandler) RequestOtpRegistHandler(c *fiber.Ctx) error {
	var req LoginorRegistRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if errs, ok := req.ValidateLoginorRegistRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta": fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Input tidak valid",
			},
			"errors": errs,
		})
	}

	_, err := h.service.SendRegistrationOTP(c.Context(), &req)
	if err != nil {
		return utils.Forbidden(c, err.Error())
	}

	return utils.Success(c, "OTP sent successfully")
}

func (h *AuthenticationHandler) VerifyOtpRegistHandler(c *fiber.Ctx) error {
	var req VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if errs, ok := req.ValidateVerifyOtpRequest(); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta":   fiber.Map{"status": fiber.StatusBadRequest, "message": "Validation error"},
			"errors": errs,
		})
	}

	stepResp, err := h.service.VerifyRegistrationOTP(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "OTP verified successfully", stepResp)
}

func (h *AuthenticationHandler) LogoutAuthentication(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	err = h.service.LogoutAuthentication(c.Context(), claims.UserID, claims.DeviceID)
	if err != nil {

		return utils.InternalServerError(c, "Failed to logout")
	}

	return utils.Success(c, "Logout successful")
}
