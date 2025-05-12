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

// func (h *RequestPickupHandler) GetAutomaticRequestByUser(c *fiber.Ctx) error {

// 	collectorId, ok := c.Locals("userID").(string)
// 	if !ok || collectorId == "" {
// 		return utils.ErrorResponse(c, "Unauthorized: User session not found")
// 	}

// 	requestPickups, err := h.service.GetAllAutomaticRequestPickup(collectorId)
// 	if err != nil {

// 		return utils.ErrorResponse(c, err.Error())
// 	}

// 	return utils.SuccessResponse(c, requestPickups, "Request pickups fetched successfully")
// }

func (h *RequestPickupHandler) GetRequestPickups(c *fiber.Ctx) error {
	// Get userID from Locals
	collectorId := c.Locals("userID").(string)

	// Call service layer to get the request pickups
	requests, err := h.service.GetRequestPickupsForCollector(collectorId)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	// Return response
	return utils.SuccessResponse(c, requests, "Automatic request pickups retrieved successfully")
}