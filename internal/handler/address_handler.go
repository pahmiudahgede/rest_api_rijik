package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type AddressHandler struct {
	AddressService services.AddressService
}

func NewAddressHandler(addressService services.AddressService) *AddressHandler {
	return &AddressHandler{AddressService: addressService}
}

func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	var requestAddressDTO dto.CreateAddressDTO
	if err := c.BodyParser(&requestAddressDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := requestAddressDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	addressResponse, err := h.AddressService.CreateAddress(c.Locals("userID").(string), requestAddressDTO)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.CreateResponse(c, addressResponse, "user address created successfully")
}

func (h *AddressHandler) GetAddressByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	addresses, err := h.AddressService.GetAddressByUserID(userID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, addresses, "User addresses fetched successfully")
}

func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	addressID := c.Params("address_id")

	address, err := h.AddressService.GetAddressByID(addressID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, address, "Address fetched successfully")
}
