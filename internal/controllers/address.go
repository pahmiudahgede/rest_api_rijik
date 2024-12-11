package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func CreateAddress(c *fiber.Ctx) error {
	var input dto.AddressInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Mohon masukkan alamat dengan benar",
			nil,
		))
	}

	if err := input.ValidatePost(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	userID := c.Locals("userID").(string)
	address, err := services.CreateAddress(userID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create address",
			nil,
		))
	}

	createdAtFormatted := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAtFormatted := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressResponse := dto.AddressResponse{
		ID:          address.ID,
		Province:    address.Province,
		District:    address.District,
		Subdistrict: address.Subdistrict,
		PostalCode:  address.PostalCode,
		Village:     address.Village,
		Detail:      address.Detail,
		Geography:   address.Geography,
		CreatedAt:   createdAtFormatted,
		UpdatedAt:   updatedAtFormatted,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Address created successfully",
		addressResponse,
	))
}

func GetListAddress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	addresses, err := services.GetAllAddressesByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Addresses not found",
			nil,
		))
	}

	var addressResponses []dto.AddressResponse
	for _, address := range addresses {

		createdAtFormatted := utils.FormatDateToIndonesianFormat(address.CreatedAt)
		updatedAtFormatted := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

		addressResponse := dto.AddressResponse{
			ID:          address.ID,
			Province:    address.Province,
			District:    address.District,
			Subdistrict: address.Subdistrict,
			PostalCode:  address.PostalCode,
			Village:     address.Village,
			Detail:      address.Detail,
			Geography:   address.Geography,
			CreatedAt:   createdAtFormatted,
			UpdatedAt:   updatedAtFormatted,
		}
		addressResponses = append(addressResponses, addressResponse)
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Addresses fetched successfully",
		addressResponses,
	))
}

func GetAddressByID(c *fiber.Ctx) error {
	addressID := c.Params("id")

	address, err := services.GetAddressByID(addressID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Address not found",
			nil,
		))
	}

	createdAtFormatted := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAtFormatted := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressResponse := dto.AddressResponse{
		ID:          address.ID,
		Province:    address.Province,
		District:    address.District,
		Subdistrict: address.Subdistrict,
		PostalCode:  address.PostalCode,
		Village:     address.Village,
		Detail:      address.Detail,
		Geography:   address.Geography,
		CreatedAt:   createdAtFormatted,
		UpdatedAt:   updatedAtFormatted,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Address fetched successfully",
		addressResponse,
	))
}

func UpdateAddress(c *fiber.Ctx) error {
	addressID := c.Params("id")

	var input dto.AddressInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
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

	address, err := services.UpdateAddress(addressID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update address",
			nil,
		))
	}

	createdAtFormatted := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAtFormatted := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressResponse := dto.AddressResponse{
		ID:          address.ID,
		Province:    address.Province,
		District:    address.District,
		Subdistrict: address.Subdistrict,
		PostalCode:  address.PostalCode,
		Village:     address.Village,
		Detail:      address.Detail,
		Geography:   address.Geography,
		CreatedAt:   createdAtFormatted,
		UpdatedAt:   updatedAtFormatted,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Address updated successfully",
		addressResponse,
	))
}

func DeleteAddress(c *fiber.Ctx) error {
	addressID := c.Params("id")

	err := services.DeleteAddress(addressID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete address",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Address deleted successfully",
		nil,
	))
}
