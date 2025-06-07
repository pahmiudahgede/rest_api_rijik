package handler
/* 
import (
	"log"
	dto "rijig/dto/auth"
	services "rijig/internal/services/auth"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthAdminHandler struct {
	UserService services.AuthAdminService
}

func NewAuthAdminHandler(userService services.AuthAdminService) *AuthAdminHandler {
	return &AuthAdminHandler{UserService: userService}
}

func (h *AuthAdminHandler) RegisterAdmin(c *fiber.Ctx) error {
	var request dto.RegisterAdminRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to parse request body")
	}

	errors, valid := request.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	user, err := h.UserService.RegisterAdmin(&request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, user, "Admin registered successfully")
}

func (h *AuthAdminHandler) LoginAdmin(c *fiber.Ctx) error {
	var request dto.LoginAdminRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to parse request body")
	}

	loginResponse, err := h.UserService.LoginAdmin(&request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, loginResponse, "Login successful")
}

func (h *AuthAdminHandler) LogoutAdmin(c *fiber.Ctx) error {
    // Ambil userID dari c.Locals
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        log.Println("Error: UserID is nil or empty")
        return utils.GenericResponse(c, fiber.StatusUnauthorized, "User not authenticated")
    }

    // Ambil deviceID dari header atau c.Locals
    deviceID, ok := c.Locals("device_id").(string)
    if !ok || deviceID == "" {
        log.Println("Error: DeviceID is nil or empty")
        return utils.ErrorResponse(c, "DeviceID is required")
    }

    log.Printf("UserID: %s, DeviceID: %s", userID, deviceID)

    err := h.UserService.LogoutAdmin(userID, deviceID)
    if err != nil {
        log.Printf("Error during logout process for user %s: %v", userID, err)
        return utils.ErrorResponse(c, err.Error())
    }

    return utils.GenericResponse(c, fiber.StatusOK, "Successfully logged out")
}

 */