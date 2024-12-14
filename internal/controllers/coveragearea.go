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
		coverageAreaResponses = append(coverageAreaResponses, dto.NewCoverageAreaResponse(
			area.ID,
			area.Province,
			utils.FormatDateToIndonesianFormat(area.CreatedAt),
			utils.FormatDateToIndonesianFormat(area.UpdatedAt),
		))
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
	coverageAreaResponse.CreatedAt = utils.FormatDateToIndonesianFormat(coverageArea.CreatedAt)
	coverageAreaResponse.UpdatedAt = utils.FormatDateToIndonesianFormat(coverageArea.UpdatedAt)

	districts, err := services.GetCoverageDistricsByCoverageAreaID(coverageArea.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch coverage districts",
			nil,
		))
	}

	var coverageAreas []dto.CoverageAreaResponse
	for _, district := range districts {
		coverageAreas = append(coverageAreas, dto.NewCoverageAreaResponse(
			district.ID,
			district.District,
			utils.FormatDateToIndonesianFormat(district.CreatedAt),
			utils.FormatDateToIndonesianFormat(district.UpdatedAt),
		))
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

	coverageArea, err := services.GetCoverageAreaByID(coverageDetail.CoverageAreaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch coverage area details by province",
			nil,
		))
	}

	subdistricts, err := services.GetSubdistrictsByCoverageDistrictID(coverageDetail.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch subdistricts",
			nil,
		))
	}

	var subdistrictResponses []dto.SubdistrictResponse
	for _, loc := range subdistricts {
		subdistrictResponses = append(subdistrictResponses, dto.NewSubdistrictResponse(
			loc.ID,
			loc.Subdistrict,
			utils.FormatDateToIndonesianFormat(loc.CreatedAt),
			utils.FormatDateToIndonesianFormat(loc.UpdatedAt),
		))
	}

	coverageAreaResponse := dto.NewCoverageAreaDetailWithLocation(
		coverageDetail.ID,
		coverageArea.Province,
		coverageDetail.District,
		utils.FormatDateToIndonesianFormat(coverageDetail.CreatedAt),
		utils.FormatDateToIndonesianFormat(coverageDetail.UpdatedAt),
		subdistrictResponses,
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage areas detail by district has been fetched",
		coverageAreaResponse,
	))
}

func CreateCoverageArea(c *fiber.Ctx) error {
	var request dto.CoverageAreaCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request payload",
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

	coverageAreaResponse := dto.NewCoverageAreaResponse(
		coverageArea.ID,
		coverageArea.Province,
		utils.FormatDateToIndonesianFormat(coverageArea.CreatedAt),
		utils.FormatDateToIndonesianFormat(coverageArea.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage area has been created",
		coverageAreaResponse,
	))
}

func CreateCoverageDistrict(c *fiber.Ctx) error {
	var request dto.CoverageDistrictCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request payload",
			nil,
		))
	}

	coverageDistrict, err := services.CreateCoverageDistrict(request.CoverageAreaID, request.District)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create coverage district",
			nil,
		))
	}

	coverageDistrictResponse := dto.NewCoverageAreaResponse(
		coverageDistrict.ID,
		coverageDistrict.District,
		utils.FormatDateToIndonesianFormat(coverageDistrict.CreatedAt),
		utils.FormatDateToIndonesianFormat(coverageDistrict.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage district has been created",
		coverageDistrictResponse,
	))
}

func CreateCoverageSubdistrict(c *fiber.Ctx) error {

	var request dto.CoverageSubdistrictCreateRequest
	if err := c.BodyParser(&request); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request payload",
			nil,
		))
	}

	coverageSubdistrict, err := services.CreateCoverageSubdistrict(
		request.CoverageAreaID,
		request.CoverageDistrictId,
		request.Subdistrict,
	)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create coverage subdistrict",
			nil,
		))
	}

	coverageSubdistrictResponse := dto.NewCoverageAreaResponse(
		coverageSubdistrict.ID,
		coverageSubdistrict.Subdistrict,
		utils.FormatDateToIndonesianFormat(coverageSubdistrict.CreatedAt),
		utils.FormatDateToIndonesianFormat(coverageSubdistrict.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Coverage subdistrict has been created",
		coverageSubdistrictResponse,
	))
}
