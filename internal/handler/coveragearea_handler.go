package handler

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CoverageAreaHandler struct {
	service services.CoverageAreaService
}

func NewCoverageAreaHandler(service services.CoverageAreaService) *CoverageAreaHandler {
	return &CoverageAreaHandler{service: service}
}

func (h *CoverageAreaHandler) CreateCoverageArea(c *fiber.Ctx) error {
	var request dto.RequestCoverageArea
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"Invalid request body"},
		})
	}

	errors, valid := request.ValidateCoverageArea()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	response, err := h.service.CreateCoverageArea(request)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Sprintf("Error creating coverage area: %v", err))
	}

	return utils.SuccessResponse(c, response, "Coverage area created successfully")
}

func (h *CoverageAreaHandler) GetCoverageAreaByID(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := h.service.GetCoverageAreaByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Coverage area with ID %s not found", id))
	}

	return utils.SuccessResponse(c, response, "Coverage area found")
}

func (h *CoverageAreaHandler) GetAllCoverageAreas(c *fiber.Ctx) error {

	response, err := h.service.GetAllCoverageAreas()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Error fetching coverage areas")
	}

	return utils.SuccessResponse(c, response, "Coverage areas fetched successfully")
}

func (h *CoverageAreaHandler) UpdateCoverageArea(c *fiber.Ctx) error {
	id := c.Params("id")
	var request dto.RequestCoverageArea
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"Invalid request body"},
		})
	}

	errors, valid := request.ValidateCoverageArea()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	response, err := h.service.UpdateCoverageArea(id, request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Coverage area with ID %s not found", id))
	}

	return utils.SuccessResponse(c, response, "Coverage area updated successfully")
}

func (h *CoverageAreaHandler) DeleteCoverageArea(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteCoverageArea(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Coverage area with ID %s not found", id))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Coverage area deleted successfully")
}
