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
