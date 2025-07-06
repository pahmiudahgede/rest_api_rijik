package wilayahindo

import (
	"strconv"
	"strings"

	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type WilayahIndonesiaHandler struct {
	WilayahService WilayahIndonesiaService
}

func NewWilayahIndonesiaHandler(wilayahService WilayahIndonesiaService) *WilayahIndonesiaHandler {
	return &WilayahIndonesiaHandler{
		WilayahService: wilayahService,
	}
}

func (h *WilayahIndonesiaHandler) ImportDataFromCSV(c *fiber.Ctx) error {
	ctx := c.Context()

	if err := h.WilayahService.ImportDataFromCSV(ctx); err != nil {
		return utils.InternalServerError(c, "Failed to import data from CSV: "+err.Error())
	}

	return utils.Success(c, "Data imported successfully from CSV")
}

func (h *WilayahIndonesiaHandler) GetAllProvinces(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	provinces, total, err := h.WilayahService.GetAllProvinces(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch provinces: "+err.Error())
	}

	response := map[string]interface{}{
		"provinces": provinces,
		"total":     total,
	}

	return utils.SuccessWithPagination(c, "Provinces retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetProvinceByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Province ID is required")
	}

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	province, totalRegencies, err := h.WilayahService.GetProvinceByID(ctx, id, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Province not found")
		}
		return utils.InternalServerError(c, "Failed to fetch province: "+err.Error())
	}

	response := map[string]interface{}{
		"province":        province,
		"total_regencies": totalRegencies,
	}

	return utils.SuccessWithPagination(c, "Province retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetAllRegencies(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	regencies, total, err := h.WilayahService.GetAllRegencies(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch regencies: "+err.Error())
	}

	response := map[string]interface{}{
		"regencies": regencies,
		"total":     total,
	}

	return utils.SuccessWithPagination(c, "Regencies retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetRegencyByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Regency ID is required")
	}

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	regency, totalDistricts, err := h.WilayahService.GetRegencyByID(ctx, id, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Regency not found")
		}
		return utils.InternalServerError(c, "Failed to fetch regency: "+err.Error())
	}

	response := map[string]interface{}{
		"regency":         regency,
		"total_districts": totalDistricts,
	}

	return utils.SuccessWithPagination(c, "Regency retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetAllDistricts(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	districts, total, err := h.WilayahService.GetAllDistricts(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch districts: "+err.Error())
	}

	response := map[string]interface{}{
		"districts": districts,
		"total":     total,
	}

	return utils.SuccessWithPagination(c, "Districts retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetDistrictByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "District ID is required")
	}

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	district, totalVillages, err := h.WilayahService.GetDistrictByID(ctx, id, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "District not found")
		}
		return utils.InternalServerError(c, "Failed to fetch district: "+err.Error())
	}

	response := map[string]interface{}{
		"district":       district,
		"total_villages": totalVillages,
	}

	return utils.SuccessWithPagination(c, "District retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetAllVillages(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	villages, total, err := h.WilayahService.GetAllVillages(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch villages: "+err.Error())
	}

	response := map[string]interface{}{
		"villages": villages,
		"total":    total,
	}

	return utils.SuccessWithPagination(c, "Villages retrieved successfully", response, page, limit)
}

func (h *WilayahIndonesiaHandler) GetVillageByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Village ID is required")
	}

	village, err := h.WilayahService.GetVillageByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Village not found")
		}
		return utils.InternalServerError(c, "Failed to fetch village: "+err.Error())
	}

	return utils.SuccessWithData(c, "Village retrieved successfully", village)
}

func (h *WilayahIndonesiaHandler) parsePaginationParams(c *fiber.Ctx) (int, int, error) {

	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid page parameter")
		}
		if parsedPage < 1 {
			return 0, 0, fiber.NewError(fiber.StatusBadRequest, "Page must be greater than 0")
		}
		page = parsedPage
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid limit parameter")
		}
		if parsedLimit < 1 {
			return 0, 0, fiber.NewError(fiber.StatusBadRequest, "Limit must be greater than 0")
		}
		if parsedLimit > 100 {
			return 0, 0, fiber.NewError(fiber.StatusBadRequest, "Limit cannot exceed 100")
		}
		limit = parsedLimit
	}

	return page, limit, nil
}

func (h *WilayahIndonesiaHandler) SetupRoutes(app *fiber.App) {

	api := app.Group("/api/v1/wilayah")

	api.Post("/import", h.ImportDataFromCSV)

	api.Get("/provinces", h.GetAllProvinces)
	api.Get("/provinces/:id", h.GetProvinceByID)

	api.Get("/regencies", h.GetAllRegencies)
	api.Get("/regencies/:id", h.GetRegencyByID)

	api.Get("/districts", h.GetAllDistricts)
	api.Get("/districts/:id", h.GetDistrictByID)

	api.Get("/villages", h.GetAllVillages)
	api.Get("/villages/:id", h.GetVillageByID)
}
