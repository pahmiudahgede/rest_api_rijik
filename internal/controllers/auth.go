package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func Register(c *fiber.Ctx) error {
	var userInput dto.RegisterUserInput

	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := userInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	err := services.RegisterUser(userInput.Username, userInput.Name, userInput.Email, userInput.Phone, userInput.Password, userInput.RoleId)
	if err != nil {

		if err.Error() == "email is already registered" {
			return c.Status(fiber.StatusConflict).JSON(utils.FormatResponse(
				fiber.StatusConflict,
				"Email is already registered",
				nil,
			))
		}
		if err.Error() == "username is already registered" {
			return c.Status(fiber.StatusConflict).JSON(utils.FormatResponse(
				fiber.StatusConflict,
				"Username is already registered",
				nil,
			))
		}
		if err.Error() == "phone number is already registered" {
			return c.Status(fiber.StatusConflict).JSON(utils.FormatResponse(
				fiber.StatusConflict,
				"Phone number is already registered",
				nil,
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create user",
			nil,
		))
	}

	user, err := repositories.GetUserByEmailOrUsername(userInput.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch user after registration",
			nil,
		))
	}

	userResponse := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"name":      user.Name,
		"email":     user.Email,
		"phone":     user.Phone,
		"roleId":    user.RoleID,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"User registered successfully",
		userResponse,
	))
}

func Login(c *fiber.Ctx) error {
	var credentials struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	token, err := services.LoginUser(credentials.EmailOrUsername, credentials.Password)
	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			err.Error(),
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Login successful",
		map[string]string{"token": token},
	))
}

func GetUserInfo(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	user, err := services.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"user tidak ditemukan",
			nil,
		))
	}

	userResponse := map[string]interface{}{
		"id":               user.ID,
		"username":         user.Username,
		"nama":             user.Name,
		"nohp":             user.Phone,
		"email":            user.Email,
		"statusverifikasi": user.EmailVerified,
		"role":             user.Role.RoleName,
		"createdAt":        user.CreatedAt,
		"updatedAt":        user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"data user berhasil ditampilkan",
		userResponse,
	))
}
