package handler

import (
	"context"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestPickupHandler interface {
	CreateRequestPickup(c *fiber.Ctx) error
	SelectCollector(c *fiber.Ctx) error
	GetAssignedPickup(c *fiber.Ctx) error
	ConfirmPickup(c *fiber.Ctx) error
	UpdatePickupStatus(c *fiber.Ctx) error
	UpdatePickupItemActualAmount(c *fiber.Ctx) error
}

type requestPickupHandler struct {
	service services.RequestPickupService
}

func NewRequestPickupHandler(service services.RequestPickupService) RequestPickupHandler {
	return &requestPickupHandler{service: service}
}

func (h *requestPickupHandler) CreateRequestPickup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req dto.RequestPickupDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if errs, ok := req.Validate(); !ok {
		return utils.ValidationErrorResponse(c, errs)
	}

	if err := h.service.ConvertCartToRequestPickup(context.Background(), userID, req); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Request pickup berhasil dibuat")
}

func (h *requestPickupHandler) SelectCollector(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID harus disertakan"},
		})
	}

	var req dto.SelectCollectorDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if errs, ok := req.Validate(); !ok {
		return utils.ValidationErrorResponse(c, errs)
	}

	if err := h.service.AssignCollectorToRequest(context.Background(), pickupID, req); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Collector berhasil dipilih untuk pickup")
}

func (h *requestPickupHandler) GetAssignedPickup(c *fiber.Ctx) error {
	collectorID := c.Locals("userID").(string)
	result, err := h.service.FindRequestsAssignedToCollector(context.Background(), collectorID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, result, "Data pickup yang ditugaskan berhasil diambil")
}

func (h *requestPickupHandler) ConfirmPickup(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID wajib diisi"},
		})
	}

	err := h.service.ConfirmPickupByCollector(context.Background(), pickupID, time.Now())
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, nil, "Pickup berhasil dikonfirmasi oleh collector")
}

func (h *requestPickupHandler) UpdatePickupStatus(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	if err := h.service.UpdatePickupStatusToPickingUp(context.Background(), pickupID); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Status pickup berhasil diperbarui menjadi 'collector_are_picking_up'")
}

func (h *requestPickupHandler) UpdatePickupItemActualAmount(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	var req dto.UpdatePickupItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if len(req.Items) == 0 {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"items": {"daftar item tidak boleh kosong"},
		})
	}

	for _, item := range req.Items {
		if item.ItemID == "" || item.Amount <= 0 {
			return utils.ValidationErrorResponse(c, map[string][]string{
				"item": {"item_id harus valid dan amount > 0"},
			})
		}
	}

	if err := h.service.UpdateActualPickupItems(context.Background(), pickupID, req.Items); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Berat aktual dan harga berhasil diperbarui")
}