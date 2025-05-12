package handler

import (
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
