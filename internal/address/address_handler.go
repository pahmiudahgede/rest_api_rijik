package address

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AddressHandler struct {
	AddressService AddressService
}

func NewAddressHandler(addressService AddressService) *AddressHandler {
	return &AddressHandler{AddressService: addressService}
}

func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	var request CreateAddressDTO
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateAddress()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}
	addressResponse, err := h.AddressService.CreateAddress(c.Context(), claims.UserID, request)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.CreateSuccessWithData(c, "user address created successfully", addressResponse)
}

func (h *AddressHandler) GetAddressByUserID(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	addresses, err := h.AddressService.GetAddressByUserID(c.Context(), claims.UserID)
	if err != nil {
		return utils.NotFound(c, err.Error())
	}

	return utils.SuccessWithData(c, "User addresses fetched successfully", addresses)
}

func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	addressID := c.Params("address_id")

	address, err := h.AddressService.GetAddressByID(c.Context(), claims.UserID, addressID)
	if err != nil {
		return utils.NotFound(c, err.Error())
	}

	return utils.SuccessWithData(c, "Address fetched successfully", address)
}

func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {

	addressID := c.Params("address_id")

	var request CreateAddressDTO
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateAddress()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	updatedAddress, err := h.AddressService.UpdateAddress(c.Context(), claims.UserID, addressID, request)
	if err != nil {
		return utils.NotFound(c, err.Error())
	}

	return utils.SuccessWithData(c, "User address updated successfully", updatedAddress)
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	addressID := c.Params("address_id")

	err = h.AddressService.DeleteAddress(c.Context(), claims.UserID, addressID)
	if err != nil {
		return utils.Forbidden(c, err.Error())
	}

	return utils.Success(c, "Address deleted successfully")
}
