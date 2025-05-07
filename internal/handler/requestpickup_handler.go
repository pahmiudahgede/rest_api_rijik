package handler

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type RequestPickupHandler struct {
	service services.RequestPickupService
}

func NewRequestPickupHandler(service services.RequestPickupService) *RequestPickupHandler {
	return &RequestPickupHandler{service: service}
}

func (h *RequestPickupHandler) CreateRequestPickup(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	var request dto.RequestPickup

	if err := c.BodyParser(&request); err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	errors, valid := request.ValidateRequestPickup()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	response, err := h.service.CreateRequestPickup(request, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Sprintf("Error creating request pickup: %v", err))
	}

	return utils.SuccessResponse(c, response, "Request pickup created successfully")
}

func (h *RequestPickupHandler) GetRequestPickupByID(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := h.service.GetRequestPickupByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Request pickup with ID %s not found: %v", id, err))
	}

	return utils.SuccessResponse(c, response, "Request pickup retrieved successfully")
}

func (h *RequestPickupHandler) GetAllRequestPickups(c *fiber.Ctx) error {

	response, err := h.service.GetAllRequestPickups()
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Sprintf("Error fetching all request pickups: %v", err))
	}

	return utils.SuccessResponse(c, response, "All request pickups retrieved successfully")
}

func (h *RequestPickupHandler) UpdateRequestPickup(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	id := c.Params("id")
	var request dto.RequestPickup

	if err := c.BodyParser(&request); err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	errors, valid := request.ValidateRequestPickup()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	response, err := h.service.UpdateRequestPickup(id, request)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Sprintf("Error updating request pickup: %v", err))
	}

	return utils.SuccessResponse(c, response, "Request pickup updated successfully")
}

func (h *RequestPickupHandler) DeleteRequestPickup(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteRequestPickup(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Request pickup with ID %s not found: %v", id, err))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Request pickup deleted successfully")
}
