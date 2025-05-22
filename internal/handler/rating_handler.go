package handler

import (
	"context"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type PickupRatingHandler interface {
	CreateRating(c *fiber.Ctx) error
	GetRatingsByCollector(c *fiber.Ctx) error
	GetAverageRating(c *fiber.Ctx) error
}

type pickupRatingHandler struct {
	service services.PickupRatingService
}

func NewPickupRatingHandler(service services.PickupRatingService) PickupRatingHandler {
	return &pickupRatingHandler{service: service}
}

func (h *pickupRatingHandler) CreateRating(c *fiber.Ctx) error {
	pickupID := c.Params("id")
	userID := c.Locals("userID").(string)
	collectorID := c.Query("collector_id")

	var req dto.CreatePickupRatingDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"Format JSON tidak valid"},
		})
	}

	if errs, ok := req.ValidateCreatePickupRatingDTO(); !ok {
		return utils.ValidationErrorResponse(c, errs)
	}

	err := h.service.CreateRating(context.Background(), userID, pickupID, collectorID, req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Rating berhasil dikirim")
}

func (h *pickupRatingHandler) GetRatingsByCollector(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	ratings, err := h.service.GetRatingsByCollector(context.Background(), collectorID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, ratings, "Daftar rating collector berhasil diambil")
}

func (h *pickupRatingHandler) GetAverageRating(c *fiber.Ctx) error {
	collectorID := c.Params("id")
	avg, err := h.service.GetAverageRating(context.Background(), collectorID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, fiber.Map{"average_rating": avg}, "Rata-rata rating collector")
}
