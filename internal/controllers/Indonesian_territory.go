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

func GetProvinceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	province, err := services.GetProvinceByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve province",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Province by id retrieved successfully",
		province,
	))
}

func GetRegencyByID(c *fiber.Ctx) error {
	id := c.Params("id")
	regency, err := services.GetRegencyByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve regency",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Regency by id retrieved successfully",
		regency,
	))
}

func GetDistrictByID(c *fiber.Ctx) error {
	id := c.Params("id")
	district, err := services.GetDistrictByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve district",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"District by id retrieved successfully",
		district,
	))
}

func GetVillageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	village, err := services.GetVillageByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to retrieve village",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Village by id retrieved successfully",
		village,
	))
}
