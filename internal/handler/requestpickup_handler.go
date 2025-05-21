package handler

// import (
// 	"fmt"
// 	"rijig/dto"
// 	"rijig/internal/services"
// 	"rijig/utils"

// 	"github.com/gofiber/fiber/v2"
// )

// type RequestPickupHandler struct {
// 	service services.RequestPickupService
// }

// func NewRequestPickupHandler(service services.RequestPickupService) *RequestPickupHandler {
// 	return &RequestPickupHandler{service: service}
// }

// func (h *RequestPickupHandler) CreateRequestPickup(c *fiber.Ctx) error {
// 	userID, ok := c.Locals("userID").(string)
// 	if !ok || userID == "" {
// 		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
// 	}

// 	var request dto.RequestPickup

// 	if err := c.BodyParser(&request); err != nil {
// 		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
// 	}

// 	errors, valid := request.ValidateRequestPickup()
// 	if !valid {
// 		return utils.ValidationErrorResponse(c, errors)
// 	}

// 	response, err := h.service.CreateRequestPickup(request, userID)
// 	if err != nil {
// 		return utils.InternalServerErrorResponse(c, fmt.Sprintf("Error creating request pickup: %v", err))
// 	}

// 	return utils.SuccessResponse(c, response, "Request pickup created successfully")
// }

// func (h *RequestPickupHandler) GetRequestPickupByID(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	response, err := h.service.GetRequestPickupByID(id)
// 	if err != nil {
// 		return utils.GenericResponse(c, fiber.StatusNotFound, fmt.Sprintf("Request pickup with ID %s not found: %v", id, err))
// 	}

// 	return utils.SuccessResponse(c, response, "Request pickup retrieved successfully")
// }

// func (h *RequestPickupHandler) GetRequestPickups(c *fiber.Ctx) error {

// 	collectorId := c.Locals("userID").(string)

// 	requests, err := h.service.GetRequestPickupsForCollector(collectorId)
// 	if err != nil {
// 		return utils.ErrorResponse(c, err.Error())
// 	}

// 	return utils.SuccessResponse(c, requests, "Automatic request pickups retrieved successfully")
// }

// func (h *RequestPickupHandler) AssignCollectorToRequest(c *fiber.Ctx) error {
// 	userId, ok := c.Locals("userID").(string)
// 	if !ok || userId == "" {
// 		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
// 	}

// 	var request dto.SelectCollectorRequest
// 	errors, valid := request.ValidateSelectCollectorRequest()
// 	if !valid {
// 		return utils.ValidationErrorResponse(c, errors)
// 	}

// 	if err := c.BodyParser(&request); err != nil {
// 		return fmt.Errorf("error parsing request body: %v", err)
// 	}

// 	err := h.service.SelectCollectorInRequest(userId, request.Collector_id)
// 	if err != nil {

// 		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error assigning collector: %v", err))
// 	}

// 	return utils.GenericResponse(c, fiber.StatusOK, "berhasil memilih collector")
// }

