package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetListUsers(c *fiber.Ctx) error {
	users, err := services.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch users",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Users fetched successfully",
		users,
	))
}

func GetUsersByRole(c *fiber.Ctx) error {
	roleID := c.Params("roleID")

	users, err := services.GetUsersByRole(roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch users by role",
			nil,
		))
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"No users found for the specified role",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Users fetched successfully",
		users,
	))
}

func GetUserByUserID(c *fiber.Ctx) error {
	userID := c.Params("userID")

	user, err := services.GetUserByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"User not found",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"User fetched successfully",
		struct {
			User dto.UserResponseDTO `json:"user"`
		}{
			User: user,
		},
	))
}
