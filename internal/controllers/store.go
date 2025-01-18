package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetStoreByID(c *fiber.Ctx) error {
	storeID := c.Params("storeid")
	if storeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Store ID is required",
			nil,
		))
	}

	store, err := services.GetStoreByID(storeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Store not found",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Store fetched successfully",
		store,
	))
}

func GetStoresByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: user ID is missing",
			nil,
		))
	}

	limit := c.QueryInt("limit", 0)
	page := c.QueryInt("page", 1)

	if limit < 0 || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid pagination parameters",
			nil,
		))
	}

	stores, err := services.GetStoresByUserID(userID, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch stores",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Stores fetched successfully",
		stores,
	))
}
