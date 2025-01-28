package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
	"golang.org/x/crypto/bcrypt"
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

	user, err := h.UserService.Login(loginDTO)
	if err != nil {
		if err.Error() == "user not found" {

			return utils.ErrorResponse(c, "User not found")
		}
		if err == bcrypt.ErrMismatchedHashAndPassword {

			return utils.ErrorResponse(c, "Invalid password")
		}

		return utils.InternalServerErrorResponse(c, "Error logging in")
	}

	return utils.LogResponse(c, user, "Login successful")
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var registerDTO dto.RegisterDTO
	if err := c.BodyParser(&registerDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := registerDTO.Validate()

	if !valid {

		return utils.ValidationErrorResponse(c, errors)
	}

	user, err := h.UserService.Register(registerDTO)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	createdAt, err := utils.FormatDateToIndonesianFormat(user.CreatedAt)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Error formatting created date")
	}

	updatedAt, err := utils.FormatDateToIndonesianFormat(user.UpdatedAt)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Error formatting updated date")
	}

	userResponse := dto.UserResponseDTO{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return utils.LogResponse(c, userResponse, "Registration successful")
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	if token == "" {

		return utils.ErrorResponse(c, "No token provided")
	}

	err := utils.DeleteData(token)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Error logging out")
	}

	return utils.NonPaginatedResponse(c, nil, 0, "Logout successful")
}
