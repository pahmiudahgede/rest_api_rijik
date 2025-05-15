package handler

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CollectorHandler struct {
	service services.CollectorService
}

func NewCollectorHandler(service services.CollectorService) *CollectorHandler {
	return &CollectorHandler{service}
}

func (h *CollectorHandler) ConfirmRequestPickup(c *fiber.Ctx) error {

	collectorId, ok := c.Locals("userID").(string)
	if !ok || collectorId == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	requestPickupId := c.Params("id")
	if requestPickupId == "" {
		return utils.ErrorResponse(c, "RequestPickup ID is required")
	}

	req, err := h.service.ConfirmRequestPickup(requestPickupId, collectorId)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, req, "Request pickup confirmed successfully")
}

func (h *CollectorHandler) GetAvaibleCollector(c *fiber.Ctx) error {

	userId := c.Locals("userID").(string)

	requests, err := h.service.FindCollectorsNearby(userId)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, requests, "menampilkan data collector terdekat")
}

func (h *CollectorHandler) ConfirmRequestManualPickup(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	requestId := c.Params("request_id")
	if requestId == "" {
		fmt.Println("requestid dibutuhkan")
	}

	var request dto.SelectCollectorRequest
	if err := c.BodyParser(&request); err != nil {
		return fmt.Errorf("error parsing request body: %v", err)
	}

	message, err := h.service.ConfirmRequestManualPickup(requestId, userId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error confirming pickup: %v", err))
	}

	return utils.SuccessResponse(c, message, "berhasil konfirmasi request")
}
