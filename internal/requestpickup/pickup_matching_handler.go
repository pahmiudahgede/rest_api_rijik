package requestpickup

import (
	"context"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type PickupMatchingHandler struct {
	service PickupMatchingService
}

func NewPickupMatchingHandler(service PickupMatchingService) *PickupMatchingHandler {
	return &PickupMatchingHandler{
		service: service,
		
	}
}

func (h *PickupMatchingHandler) GetNearbyCollectorsForPickup(c *fiber.Ctx) error {
	pickupID := c.Params("pickupID")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validasi gagal", map[string][]string{
			"pickup_id": {"pickup ID harus disertakan"},
		})
	}

	collectors, err := h.service.FindNearbyCollectorsForPickup(context.Background(), pickupID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Data collector terdekat berhasil diambil", collectors)
}

func (h *PickupMatchingHandler) GetAvailablePickupForCollector(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	pickups, err := h.service.FindAvailableRequestsForCollector(context.Background(), claims.UserID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Data request pickup otomatis berhasil diambil", pickups)
}
