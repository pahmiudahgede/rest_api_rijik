package handler
/* 
import (
	"log"
	"rijig/dto"
	services "rijig/internal/services/auth"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthMasyarakatHandler struct {
	authMasyarakatService services.AuthMasyarakatService
}

func NewAuthMasyarakatHandler(authMasyarakatService services.AuthMasyarakatService) *AuthMasyarakatHandler {
	return &AuthMasyarakatHandler{authMasyarakatService}
}

func (h *AuthMasyarakatHandler) RegisterOrLoginHandler(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	if req.Phone == "" {
		return utils.ErrorResponse(c, "Phone number is required")
	}

	if err := h.authMasyarakatService.RegisterOrLogin(&req); err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "OTP sent successfully")
}

func (h *AuthMasyarakatHandler) VerifyOTPHandler(c *fiber.Ctx) error {
	var req dto.VerifyOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, "Invalid request body")
	}

	if req.OTP == "" {
		return utils.ErrorResponse(c, "OTP is required")
	}

	if req.DeviceID == "" {
		return utils.ErrorResponse(c, "DeviceID is required")
	}

	response, err := h.authMasyarakatService.VerifyOTP(&req)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, response, "Registration/Login successful")
}

func (h *AuthMasyarakatHandler) LogoutHandler(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "User is not logged in or invalid session")
	}

	deviceID, ok := c.Locals("device_id").(string)
	if !ok || deviceID == "" {
		log.Println("Error: DeviceID is nil or empty")
		return utils.ErrorResponse(c, "DeviceID is required")
	}

	err := h.authMasyarakatService.Logout(userID, deviceID)
	if err != nil {
		log.Printf("Error during logout process for user %s: %v", userID, err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Logged out successfully")
}
 */