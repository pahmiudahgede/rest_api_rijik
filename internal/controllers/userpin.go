package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func CreatePin(c *fiber.Ctx) error {
	var input dto.PinInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Data input tidak valid",
			nil,
		))
	}

	if err := input.ValidateCreate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	userID := c.Locals("userID").(string)

	existingPin, err := services.GetPinByUserID(userID)
	if err == nil && existingPin.ID != "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"PIN sudah ada, tidak perlu dibuat lagi",
			nil,
		))
	}

	pin, err := services.CreatePin(userID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create PIN",
			nil,
		))
	}

	pinResponse := map[string]interface{}{
		"id":        pin.ID,
		"createdAt": pin.CreatedAt,
		"updatedAt": pin.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"PIN created successfully",
		pinResponse,
	))
}

func GetPin(c *fiber.Ctx) error {
	var input dto.PinInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Data input tidak valid",
			nil,
		))
	}

	userID := c.Locals("userID").(string)
	pin, err := services.GetPinByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"PIN tidak ditemukan",
			nil,
		))
	}

	isPinValid := services.CheckPin(pin.Pin, input.Pin)

	if isPinValid {
		return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
			fiber.StatusOK,
			"PIN benar",
			true,
		))
	}

	return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
		fiber.StatusUnauthorized,
		"PIN salah",
		false,
	))
}

func UpdatePin(c *fiber.Ctx) error {
	var input dto.PinUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Data input tidak valid",
			nil,
		))
	}

	if err := input.ValidateUpdate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	userID := c.Locals("userID").(string)

	updatedPin, err := services.UpdatePin(userID, input.OldPin, input.NewPin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update PIN",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"PIN updated successfully",
		updatedPin,
	))
}
