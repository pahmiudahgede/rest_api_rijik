package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetUserInitialCoint(c *fiber.Ctx) error {
	points, err := services.GetPoints()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch points",
			nil,
		))
	}

	var pointResponses []dto.PointResponse
	for _, point := range points {
		pointResponses = append(pointResponses, dto.PointResponse{
			ID:           point.ID,
			CoinName:     point.CoinName,
			ValuePerUnit: point.ValuePerUnit,
			CreatedAt:    utils.FormatDateToIndonesianFormat(point.CreatedAt),
			UpdatedAt:    utils.FormatDateToIndonesianFormat(point.UpdatedAt),
		})
	}

	if len(pointResponses) == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
			fiber.StatusOK,
			"Points successfully displayed but no data",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Points fetched successfully",
		struct {
			Points []dto.PointResponse `json:"points"`
		}{
			Points: pointResponses,
		},
	))

}

func GetUserInitialCointById(c *fiber.Ctx) error {
	id := c.Params("id")

	point, err := services.GetPointByID(id)
	if err != nil {
		if err.Error() == "point not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
				fiber.StatusNotFound,
				"Point not found",
				nil,
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch point",
			nil,
		))
	}

	pointResponse := dto.PointResponse{
		ID:           point.ID,
		CoinName:     point.CoinName,
		ValuePerUnit: point.ValuePerUnit,
		CreatedAt:    utils.FormatDateToIndonesianFormat(point.CreatedAt),
		UpdatedAt:    utils.FormatDateToIndonesianFormat(point.UpdatedAt),
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point fetched successfully",
		pointResponse,
	))
}

func CreatePoint(c *fiber.Ctx) error {
	var pointInput dto.PointRequest

	if err := c.BodyParser(&pointInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := pointInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	newPoint, err := services.CreatePoint(pointInput.CoinName, pointInput.ValuePerUnit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create point",
			nil,
		))
	}

	pointResponse := dto.NewPointResponse(
		newPoint.ID,
		newPoint.CoinName,
		newPoint.ValuePerUnit,
		utils.FormatDateToIndonesianFormat(newPoint.CreatedAt),
		utils.FormatDateToIndonesianFormat(newPoint.UpdatedAt),
	)

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Point created successfully",
		struct {
			Point dto.PointResponse `json:"point"`
		}{
			Point: pointResponse,
		},
	))
}

func UpdatePoint(c *fiber.Ctx) error {
	id := c.Params("id")

	var pointInput dto.PointUpdateDTO

	if err := c.BodyParser(&pointInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := pointInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	updatedPoint, err := services.UpdatePoint(id, pointInput.CoinName, pointInput.ValuePerUnit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update point",
			nil,
		))
	}

	pointResponse := dto.NewPointResponse(
		updatedPoint.ID,
		updatedPoint.CoinName,
		updatedPoint.ValuePerUnit,
		utils.FormatDateToIndonesianFormat(updatedPoint.CreatedAt),
		utils.FormatDateToIndonesianFormat(updatedPoint.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point updated successfully",
		struct {
			Point dto.PointResponse `json:"point"`
		}{
			Point: pointResponse,
		},
	))
}

func DeletePoint(c *fiber.Ctx) error {
	id := c.Params("id")

	err := services.DeletePoint(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete point",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Point deleted successfully",
		nil,
	))
}
