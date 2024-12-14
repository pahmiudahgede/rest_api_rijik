package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetProvinces(c *fiber.Ctx) error {
	provinces, err := services.GetProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve provinces",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Provinces retrieved successfully",
		provinces,
	))
}

func GetRegencies(c *fiber.Ctx) error {
	regencies, err := services.GetRegencies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve regencies",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Regencies retrieved successfully",
		regencies,
	))
}

func GetDistricts(c *fiber.Ctx) error {
	districts, err := services.GetDistricts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve districts",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Districts retrieved successfully",
		districts,
	))
}

func GetVillages(c *fiber.Ctx) error {
	villages, err := services.GetVillages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve villages",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Villages retrieved successfully",
		villages,
	))
}
