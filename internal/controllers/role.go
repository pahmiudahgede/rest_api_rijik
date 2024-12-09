package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetUserRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	role, err := services.GetUserRoleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"UserRole tidak ditemukan",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"UserRole ditemukan",
		role,
	))
}

func GetAllUserRoles(c *fiber.Ctx) error {
	roles, err := services.GetAllUserRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Gagal mengambil data UserRole",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Daftar UserRole",
		roles,
	))
}
