package handler

import (
	"context"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CollectorHandler interface {
	CreateCollector(c *fiber.Ctx) error
	AddTrashToCollector(c *fiber.Ctx) error
	GetCollectorByID(c *fiber.Ctx) error
	GetCollectorByUserID(c *fiber.Ctx) error
	UpdateCollector(c *fiber.Ctx) error
	UpdateJobStatus(c *fiber.Ctx) error
	UpdateTrash(c *fiber.Ctx) error
	DeleteTrash(c *fiber.Ctx) error
}
type collectorHandler struct {
	service services.CollectorService
}

func NewCollectorHandler(service services.CollectorService) CollectorHandler {
	return &collectorHandler{service: service}
}

// func (h *CollectorHandler) ConfirmRequestPickup(c *fiber.Ctx) error {

// 	collectorId, ok := c.Locals("userID").(string)
// 	if !ok || collectorId == "" {
// 		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
// 	}

// 	requestPickupId := c.Params("id")
// 	if requestPickupId == "" {
// 		return utils.ErrorResponse(c, "RequestPickup ID is required")
// 	}

// 	req, err := h.service.ConfirmRequestPickup(requestPickupId, collectorId)
// 	if err != nil {
// 		return utils.ErrorResponse(c, err.Error())
// 	}

// 	return utils.SuccessResponse(c, req, "Request pickup confirmed successfully")
// }

// func (h *CollectorHandler) GetAvaibleCollector(c *fiber.Ctx) error {

// 	userId := c.Locals("userID").(string)

// 	requests, err := h.service.FindCollectorsNearby(userId)
// 	if err != nil {
// 		return utils.ErrorResponse(c, err.Error())
// 	}

// 	return utils.SuccessResponse(c, requests, "menampilkan data collector terdekat")
// }

// func (h *CollectorHandler) ConfirmRequestManualPickup(c *fiber.Ctx) error {
// 	userId := c.Locals("userID").(string)
// 	requestId := c.Params("request_id")
// 	if requestId == "" {
// 		fmt.Println("requestid dibutuhkan")
// 	}

// 	var request dto.SelectCollectorRequest
// 	if err := c.BodyParser(&request); err != nil {
// 		return fmt.Errorf("error parsing request body: %v", err)
// 	}

// 	message, err := h.service.ConfirmRequestManualPickup(requestId, userId)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error confirming pickup: %v", err))
// 	}

// 	return utils.SuccessResponse(c, message, "berhasil konfirmasi request")
// }

func (h *collectorHandler) CreateCollector(c *fiber.Ctx) error {
	var req dto.RequestCollectorDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if errs, valid := req.ValidateRequestCollector(); !valid {
		return utils.ValidationErrorResponse(c, errs)
	}

	userID := c.Locals("userID").(string)
	err := h.service.CreateCollector(context.Background(), userID, req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.CreateResponse(c, nil, "Collector berhasil dibuat")
}

// POST /collectors/:id/trash
func (h *collectorHandler) AddTrashToCollector(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	var req dto.RequestAddAvaibleTrash

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if errs, valid := req.ValidateRequestAddAvaibleTrash(); !valid {
		return utils.ValidationErrorResponse(c, errs)
	}

	err := h.service.AddTrashToCollector(context.Background(), collectorID, req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Trash berhasil ditambahkan")
}

// GET /collectors/:id
func (h *collectorHandler) GetCollectorByID(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	result, err := h.service.GetCollectorByID(context.Background(), collectorID)
	if err != nil {
		return utils.ErrorResponse(c, "Collector tidak ditemukan")
	}
	return utils.SuccessResponse(c, result, "Data collector berhasil diambil")
}
func (h *collectorHandler) GetCollectorByUserID(c *fiber.Ctx) error {
	// collectorID := c.Params("id")
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	result, err := h.service.GetCollectorByUserID(context.Background(), userID)
	if err != nil {
		return utils.ErrorResponse(c, "Collector tidak ditemukan")
	}
	return utils.SuccessResponse(c, result, "Data collector berhasil diambil")
}

// PATCH /collectors/:id
func (h *collectorHandler) UpdateCollector(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	var req struct {
		JobStatus *string `json:"job_status"`
		Rating    float32 `json:"rating"`
		AddressID string  `json:"address_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if req.AddressID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"address_id": {"tidak boleh kosong"},
		})
	}

	err := h.service.UpdateCollector(context.Background(), collectorID, req.JobStatus, req.Rating, req.AddressID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Collector berhasil diperbarui")
}

func (h *collectorHandler) UpdateJobStatus(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	var req struct {
		JobStatus string `json:"job_status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	if req.JobStatus != "active" && req.JobStatus != "inactive" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"job_status": {"harus bernilai 'active' atau 'inactive'"},
		})
	}

	err := h.service.UpdateCollector(c.Context(), collectorID, &req.JobStatus, 0, "")
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Status collector berhasil diperbarui")
}


// PATCH /collectors/:id/trash
func (h *collectorHandler) UpdateTrash(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	var req []dto.RequestAvaibleTrashbyCollector

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format JSON tidak valid"},
		})
	}

	for i, t := range req {
		if t.TrashId == "" {
			return utils.ValidationErrorResponse(c, map[string][]string{
				"trash_id": {t.TrashId, "trash_id tidak boleh kosong pada item ke " + string(rune(i))},
			})
		}
		if t.TrashPrice <= 0 {
			return utils.ValidationErrorResponse(c, map[string][]string{
				"trash_price": {"trash_price harus lebih dari 0 pada item ke " + string(rune(i))},
			})
		}
	}

	err := h.service.UpdateAvaibleTrashByCollector(context.Background(), collectorID, req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Trash berhasil diperbarui")
}

// DELETE /collectors/trash/:id
func (h *collectorHandler) DeleteTrash(c *fiber.Ctx) error {
	trashID := c.Params("id")
	if trashID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"trash_id": {"tidak boleh kosong"},
		})
	}

	err := h.service.DeleteAvaibleTrash(context.Background(), trashID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Trash berhasil dihapus")
}
