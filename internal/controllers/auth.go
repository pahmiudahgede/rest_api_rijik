package controllers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
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

	err := services.RegisterUser(userInput.Username, userInput.Name, userInput.Email, userInput.Phone, userInput.Password, userInput.ConfirmPassword, userInput.RoleId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(utils.FormatResponse(
			fiber.StatusConflict,
			err.Error(),
			nil,
		))
	}

	user, err := repositories.GetUserByEmailUsernameOrPhone(userInput.Email, userInput.RoleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch user after registration",
			nil,
		))
	}

	userResponse := dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		RoleId:    user.RoleID,
		CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
		UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"User registered successfully",
		userResponse,
	))
}

func Login(c *fiber.Ctx) error {
	var credentials struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	token, err := services.LoginUser(credentials.Identifier, credentials.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			err.Error(),
			nil,
		))
	}

	ctx := context.Background()
	err = config.RedisClient.Set(ctx, "auth_token:"+token, credentials.Identifier, time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to store session",
			nil,
		))
	}

	user, err := repositories.GetUserByEmailUsernameOrPhone(credentials.Identifier, "")
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
		map[string]interface{}{
			"token": token,
			"role":  user.RoleID,
		},
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

	userResponse := dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		RoleId:    user.RoleID,
		CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
		UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Data user berhasil ditampilkan",
		userResponse,
	))
}

func UpdateUser(c *fiber.Ctx) error {
	var userInput dto.UpdateUserInput

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

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized access",
			nil,
		))
	}

	err := services.UpdateUser(userID, userInput.Email, userInput.Username, userInput.Name, userInput.Phone)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(utils.FormatResponse(
			fiber.StatusConflict,
			err.Error(),
			nil,
		))
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch user after update",
			nil,
		))
	}

	userResponse := dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		RoleId:    user.RoleID,
		CreatedAt: utils.FormatDateToIndonesianFormat(user.CreatedAt),
		UpdatedAt: utils.FormatDateToIndonesianFormat(user.UpdatedAt),
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"User updated successfully",
		userResponse,
	))
}

func UpdatePassword(c *fiber.Ctx) error {
	var passwordInput dto.UpdatePasswordInput

	if err := c.BodyParser(&passwordInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := passwordInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	userID := c.Locals("userID").(string)

	err := services.UpdatePassword(userID, passwordInput.OldPassword, passwordInput.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			err.Error(),
			nil,
		))
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch user after password update",
			nil,
		))
	}

	updatedAtFormatted := utils.FormatDateToIndonesianFormat(user.UpdatedAt)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Password updated successfully",
		map[string]string{
			"updatedAt": updatedAtFormatted,
		},
	))
}

func Logout(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Token is required",
			nil,
		))
	}

	ctx := context.Background()
	err := config.RedisClient.Del(ctx, "auth_token:"+tokenString).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete session",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Logout successful",
		nil,
	))
}
