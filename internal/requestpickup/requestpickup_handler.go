package requestpickup

import (
	"context"
	"rijig/middleware"
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
	service RequestPickupService
}

func NewRequestPickupHandler(service RequestPickupService) RequestPickupHandler {
	return &requestPickupHandler{service: service}
}

func (h *requestPickupHandler) CreateRequestPickup(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req RequestPickupDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := req.Validate()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	if err := h.service.ConvertCartToRequestPickup(context.Background(), claims.UserID, req); err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Request pickup berhasil dibuat")
}

func (h *requestPickupHandler) SelectCollector(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"pickup_id": {"pickup ID harus disertakan"},
		})
	}

	var req SelectCollectorDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if errs, ok := req.Validate(); !ok {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errs)
	}

	if err := h.service.AssignCollectorToRequest(context.Background(), pickupID, req); err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Collector berhasil dipilih untuk pickup")
}

func (h *requestPickupHandler) GetAssignedPickup(c *fiber.Ctx) error {
	// collectorID := c.Locals("userID").(string)
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	result, err := h.service.FindRequestsAssignedToCollector(context.Background(), claims.UserID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.SuccessWithData(c, "Data pickup yang ditugaskan berhasil diambil", result)
}

func (h *requestPickupHandler) ConfirmPickup(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"pickup_id": {"pickup ID wajib diisi"},
		})
	}

	err := h.service.ConfirmPickupByCollector(context.Background(), pickupID, time.Now())
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, "Pickup berhasil dikonfirmasi oleh collector")
}

func (h *requestPickupHandler) UpdatePickupStatus(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	if err := h.service.UpdatePickupStatusToPickingUp(context.Background(), pickupID); err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Status pickup berhasil diperbarui menjadi 'collector_are_picking_up'")
}

func (h *requestPickupHandler) UpdatePickupItemActualAmount(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	if pickupID == "" {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"pickup_id": {"pickup ID tidak boleh kosong"},
		})
	}

	var req UpdatePickupItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if len(req.Items) == 0 {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
			"items": {"daftar item tidak boleh kosong"},
		})
	}

	for _, item := range req.Items {
		if item.ItemID == "" || item.Amount <= 0 {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", map[string][]string{
				"item": {"item_id harus valid dan amount > 0"},
			})
		}
	}

	if err := h.service.UpdateActualPickupItems(context.Background(), pickupID, req.Items); err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Berat aktual dan harga berhasil diperbarui")
}
