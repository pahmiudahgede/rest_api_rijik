package handler

import (
	"context"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type PickupMatchingHandler interface {
	GetNearbyCollectorsForPickup(c *fiber.Ctx) error
	GetAvailablePickupForCollector(c *fiber.Ctx) error
}

type pickupMatchingHandler struct {
	service services.PickupMatchingService
}

func NewPickupMatchingHandler(service services.PickupMatchingService) PickupMatchingHandler {
	return &pickupMatchingHandler{service: service}
}

func (h *pickupMatchingHandler) GetNearbyCollectorsForPickup(c *fiber.Ctx) error {
	pickupID := c.Params("pickupID")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID harus disertakan"},
		})
	}

	collectors, err := h.service.FindNearbyCollectorsForPickup(context.Background(), pickupID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, collectors, "Data collector terdekat berhasil diambil")
}

func (h *pickupMatchingHandler) GetAvailablePickupForCollector(c *fiber.Ctx) error {
	collectorID := c.Locals("userID").(string)

	pickups, err := h.service.FindAvailableRequestsForCollector(context.Background(), collectorID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, pickups, "Data request pickup otomatis berhasil diambil")
}
