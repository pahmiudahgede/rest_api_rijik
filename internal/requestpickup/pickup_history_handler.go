package requestpickup

import (
	"context"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type PickupStatusHistoryHandler interface {
	GetStatusHistory(c *fiber.Ctx) error
}

type pickupStatusHistoryHandler struct {
	service PickupStatusHistoryService
}

func NewPickupStatusHistoryHandler(service PickupStatusHistoryService) PickupStatusHistoryHandler {
	return &pickupStatusHistoryHandler{service: service}
}

func (h *pickupStatusHistoryHandler) GetStatusHistory(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	histories, err := h.service.GetStatusHistory(context.Background(), pickupID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Riwayat status pickup berhasil diambil", histories)
}