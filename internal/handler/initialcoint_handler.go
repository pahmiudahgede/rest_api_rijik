package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type InitialCointHandler struct {
	InitialCointService services.InitialCointService
}

func NewInitialCointHandler(initialCointService services.InitialCointService) *InitialCointHandler {
	return &InitialCointHandler{InitialCointService: initialCointService}
}

func (h *InitialCointHandler) CreateInitialCoint(c *fiber.Ctx) error {
	var request dto.RequestInitialCointDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateCointInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	initialCointResponse, err := h.InitialCointService.CreateInitialCoint(request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.CreateResponse(c, initialCointResponse, "Initial coint created successfully")
}

func (h *InitialCointHandler) GetAllInitialCoints(c *fiber.Ctx) error {
	initialCoints, err := h.InitialCointService.GetAllInitialCoints()
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to fetch initial coints")
	}

	return utils.NonPaginatedResponse(c, initialCoints, len(initialCoints), "Initial coints fetched successfully")
}

func (h *InitialCointHandler) GetInitialCointByID(c *fiber.Ctx) error {
	id := c.Params("coin_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Coin ID is required")
	}

	initialCoint, err := h.InitialCointService.GetInitialCointByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "Invalid coin ID")
	}

	return utils.SuccessResponse(c, initialCoint, "Initial coint fetched successfully")
}

func (h *InitialCointHandler) UpdateInitialCoint(c *fiber.Ctx) error {
	id := c.Params("coin_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Coin ID is required")
	}

	var request dto.RequestInitialCointDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateCointInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	initialCointResponse, err := h.InitialCointService.UpdateInitialCoint(id, request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, initialCointResponse, "Initial coint updated successfully")
}

func (h *InitialCointHandler) DeleteInitialCoint(c *fiber.Ctx) error {
	id := c.Params("coin_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Coin ID is required")
	}

	err := h.InitialCointService.DeleteInitialCoint(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Initial coint deleted successfully")
}
