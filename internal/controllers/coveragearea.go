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
		District:         coverageDetail.District,
		LocationSpecific: locationSpecificResponses,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage areas detail by district has been fetched",
		coverageAreaResponse,
	))
}

func CreateCoverageArea(c *fiber.Ctx) error {
	var request dto.CoverageAreaRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	coverageArea, err := services.CreateCoverageArea(request.Province)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create coverage area",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage area has been created successfully",
		coverageArea,
	))
}

func CreateCoverageDetail(c *fiber.Ctx) error {
	var request dto.CoverageDetailRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	coverageDetail, err := services.CreateCoverageDetail(request.CoverageAreaID, request.Province, request.District)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create coverage detail",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage detail has been created successfully",
		coverageDetail,
	))
}

func CreateLocationSpecific(c *fiber.Ctx) error {
	var request dto.LocationSpecificRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	locationSpecific, err := services.CreateLocationSpecific(request.CoverageDetailID, request.Subdistrict)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create location specific",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Location specific has been created successfully",
		locationSpecific,
	))
}
