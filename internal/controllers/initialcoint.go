package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type PointController struct {
	service *services.PointService
}

func NewPointController(service *services.PointService) *PointController {
	return &PointController{service: service}
}

func (pc *PointController) GetAllPoints(c *fiber.Ctx) error {
	points, err := pc.service.GetAllPoints()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch points",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Points fetched successfully",
		points,
	))
}

func (pc *PointController) GetPointByID(c *fiber.Ctx) error {
	id := c.Params("id")
	point, err := pc.service.GetPointByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
			fiber.StatusNotFound,
			"Point not found",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point fetched successfully",
		point,
	))
}

func (pc *PointController) CreatePoint(c *fiber.Ctx) error {
	var request dto.PointCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	point, err := pc.service.CreatePoint(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Point created successfully",
		point,
	))
}

func (pc *PointController) UpdatePoint(c *fiber.Ctx) error {
	id := c.Params("id")
	var request dto.PointUpdateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	point, err := pc.service.UpdatePoint(id, &request)
	if err != nil {
		if err.Error() == "point not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Point not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point updated successfully",
		point,
	))
}

func (pc *PointController) DeletePoint(c *fiber.Ctx) error {
	id := c.Params("id")

	err := pc.service.DeletePoint(id)
	if err != nil {
		if err.Error() == "point not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Point not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point deleted successfully",
		nil,
	))
}