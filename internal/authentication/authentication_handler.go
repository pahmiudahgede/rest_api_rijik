package authentication

import (
	"log"
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
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	if claims.DeviceID == "" {
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

	if claims.UserID == "" {
		return utils.BadRequest(c, "userid is required")
	}

	tokenData, err := utils.RefreshAccessToken(claims.UserID, claims.DeviceID, body.RefreshToken)
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

func (h *AuthenticationHandler) GetRegistrationStatus(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "unauthorized access")
	}

	res, err := h.service.GetRegistrationStatus(c.Context(), claims.UserID, claims.DeviceID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Registration status retrieved successfully", res)
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

func (h *AuthenticationHandler) RegisterAdmin(c *fiber.Ctx) error {

	var req RegisterAdminRequest
	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	response, err := h.service.RegisterAdmin(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "Admin registered successfully", response)
}

// POST /auth/admin/verify-email - Verify email dari registration
func (h *AuthenticationHandler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	err := h.service.VerifyEmail(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "Email berhasil diverifikasi. Sekarang Anda dapat login", nil)
}

// POST /auth/admin/resend-verification - Resend verification email
func (h *AuthenticationHandler) ResendEmailVerification(c *fiber.Ctx) error {
	var req ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	response, err := h.service.ResendEmailVerification(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "Verification email resent", response)
}

func (h *AuthenticationHandler) VerifyAdminOTP(c *fiber.Ctx) error {
	var req VerifyAdminOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if errs, ok := req.Valida(); !ok {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"meta": fiber.Map{
	// 			"status":  fiber.StatusBadRequest,
	// 			"message": "periksa lagi inputan",
	// 		},
	// 		"errors": errs,
	// 	})
	// }

	response, err := h.service.VerifyAdminOTP(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "OTP resent successfully", response)
}

// POST /auth/admin/resend-otp - Resend OTP
func (h *AuthenticationHandler) ResendAdminOTP(c *fiber.Ctx) error {
	var req ResendAdminOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	response, err := h.service.ResendAdminOTP(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "OTP resent successfully", response)
}

func (h *AuthenticationHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	response, err := h.service.ForgotPassword(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "Reset password email sent", response)
}

// POST /auth/admin/reset-password - Step 2: Reset password dengan token
func (h *AuthenticationHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// if err := h.validator.Struct(&req); err != nil {
	// 	return utils.BadRequest(c, "Validation failed: "+err.Error())
	// }

	err := h.service.ResetPassword(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "Password berhasil direset", nil)
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
