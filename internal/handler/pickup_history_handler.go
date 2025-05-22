package handler

import (
	"context"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type PickupStatusHistoryHandler interface {
	GetStatusHistory(c *fiber.Ctx) error
}

type pickupStatusHistoryHandler struct {
	service services.PickupStatusHistoryService
}

func NewPickupStatusHistoryHandler(service services.PickupStatusHistoryService) PickupStatusHistoryHandler {
	return &pickupStatusHistoryHandler{service: service}
}

func (h *pickupStatusHistoryHandler) GetStatusHistory(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	histories, err := h.service.GetStatusHistory(context.Background(), pickupID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, histories, "Riwayat status pickup berhasil diambil")
}
