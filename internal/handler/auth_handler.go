package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginDTO dto.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	validationErrors, valid := loginDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	user, err := h.UserService.Login(loginDTO)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, user, "Login successful")
}

func (h *UserHandler) Register(c *fiber.Ctx) error {

	var registerDTO dto.RegisterDTO
	if err := c.BodyParser(&registerDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid request body"}})
	}

	errors, valid := registerDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userResponse, err := h.UserService.Register(registerDTO)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.CreateResponse(c, userResponse, "Registration successful")
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Println("Unauthorized access: User ID not found in session")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	err := utils.DeleteSessionData(userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Error logging out")
	}

	return utils.SuccessResponse(c, nil, "Logout successful")
}
