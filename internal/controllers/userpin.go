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

	formattedCreatedAt := utils.FormatDateToIndonesianFormat(pin.CreatedAt)

	pinResponse := dto.PinResponse{

		CreatedAt: formattedCreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"PIN created successfully",
		pinResponse,
	))
}

func GetPinStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	pin, err := services.GetPinByUserID(userID)
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Anda belum membuat PIN",
			nil,
		))
	}

	formattedCreatedAt := utils.FormatDateToIndonesianFormat(pin.CreatedAt)
	formattedUpdatedAt := utils.FormatDateToIndonesianFormat(pin.UpdatedAt)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"PIN sudah dibuat",
		map[string]interface{}{
			"createdAt": formattedCreatedAt,
			"updatedAt": formattedUpdatedAt,
		},
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
			"Sepertinya anda belum membuat pin",
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

		if err.Error() == "PIN lama salah" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
				fiber.StatusUnauthorized,
				"PIN lama salah",
				nil,
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update PIN",
			nil,
		))
	}

	formattedUpdatedAt := utils.FormatDateToIndonesianFormat(updatedPin.UpdatedAt)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"PIN updated successfully",
		map[string]interface{}{
			"id":        updatedPin.ID,
			"updatedAt": formattedUpdatedAt,
		},
	))
}
