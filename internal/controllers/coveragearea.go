package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetCoverageAreas(c *fiber.Ctx) error {
	coverageAreas, err := services.GetCoverageAreas()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch coverage areas",
			nil,
		))
	}

	var coverageAreaResponses []dto.CoverageAreaResponse
	for _, area := range coverageAreas {
		coverageAreaResponses = append(coverageAreaResponses, dto.NewCoverageAreaResponse(area.ID, area.Province))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage areas has been fetched",
		coverageAreaResponses,
	))
}

func GetCoverageAreaByIDProvince(c *fiber.Ctx) error {
	id := c.Params("id")

	coverageArea, err := services.GetCoverageAreaByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch coverage area details by province",
			nil,
		))
	}

	var coverageAreaResponse dto.CoverageAreaWithDistrictsResponse
	coverageAreaResponse.ID = coverageArea.ID
	coverageAreaResponse.Province = coverageArea.Province

	var coverageAreas []dto.CoverageAreaResponse
	for _, detail := range coverageArea.Details {
		coverageAreas = append(coverageAreas, dto.NewCoverageAreaResponse(detail.ID, detail.District))
	}

	coverageAreaResponse.CoverageArea = coverageAreas

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage areas detail by province has been fetched",
		coverageAreaResponse,
	))
}

func GetCoverageAreaByIDDistrict(c *fiber.Ctx) error {
	id := c.Params("id")

	coverageDetail, err := services.GetCoverageAreaByDistrictID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch coverage area details by district",
			nil,
		))
	}

	var locationSpecificResponses []dto.LocationSpecificResponse
	for _, loc := range coverageDetail.LocationSpecific {
		locationSpecificResponses = append(locationSpecificResponses, dto.NewLocationSpecificResponse(loc.ID, loc.Subdistrict))
	}

	coverageAreaResponse := dto.CoverageAreaDetailWithLocation{
		ID:               coverageDetail.ID,
		Province:         coverageDetail.Province,
		District:         coverageDetail.District,
		LocationSpecific: locationSpecificResponses,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage areas detail by district has been fetched",
		coverageAreaResponse,
	))
}
