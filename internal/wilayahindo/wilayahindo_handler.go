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

	page, limit, isPaginated, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	provinces, total, err := h.WilayahService.GetAllProvinces(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch provinces: "+err.Error())
	}

	if isPaginated {
		return utils.SuccessWithPaginationAndTotal(c, "Provinces retrieved successfully", provinces, page, limit, total)
	}
	return utils.SuccessWithTotal(c, "Provinces retrieved successfully", provinces, total)
}

func (h *WilayahIndonesiaHandler) GetProvinceByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Province ID is required")
	}

	page, limit, isPaginated, err := h.parsePaginationParams(c)
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

	// Add total_regencies info to meta through custom response
	responseData := map[string]interface{}{
		"id":        province.ID,
		"name":      province.Name,
		"regencies": province.Regencies,
	}

	if isPaginated {
		return h.responseWithCustomMeta(c, "Province retrieved successfully", responseData, page, limit, totalRegencies)
	}
	return h.responseWithCustomMetaNoPage(c, "Province retrieved successfully", responseData, totalRegencies)
}

func (h *WilayahIndonesiaHandler) GetAllRegencies(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, isPaginated, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	regencies, total, err := h.WilayahService.GetAllRegencies(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch regencies: "+err.Error())
	}

	if isPaginated {
		return utils.SuccessWithPaginationAndTotal(c, "Regencies retrieved successfully", regencies, page, limit, total)
	}
	return utils.SuccessWithTotal(c, "Regencies retrieved successfully", regencies, total)
}

func (h *WilayahIndonesiaHandler) GetRegencyByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Regency ID is required")
	}

	page, limit, isPaginated, err := h.parsePaginationParams(c)
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

	responseData := map[string]interface{}{
		"id":          regency.ID,
		"province_id": regency.ProvinceID,
		"name":        regency.Name,
		"districts":   regency.Districts,
	}

	if isPaginated {
		return h.responseWithCustomMeta(c, "Regency retrieved successfully", responseData, page, limit, totalDistricts)
	}
	return h.responseWithCustomMetaNoPage(c, "Regency retrieved successfully", responseData, totalDistricts)
}

func (h *WilayahIndonesiaHandler) GetAllDistricts(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, isPaginated, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	districts, total, err := h.WilayahService.GetAllDistricts(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch districts: "+err.Error())
	}

	if isPaginated {
		return utils.SuccessWithPaginationAndTotal(c, "Districts retrieved successfully", districts, page, limit, total)
	}
	return utils.SuccessWithTotal(c, "Districts retrieved successfully", districts, total)
}

func (h *WilayahIndonesiaHandler) GetDistrictByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "District ID is required")
	}

	page, limit, isPaginated, err := h.parsePaginationParams(c)
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

	responseData := map[string]interface{}{
		"id":         district.ID,
		"regency_id": district.RegencyID,
		"name":       district.Name,
		"villages":   district.Villages,
	}

	if isPaginated {
		return h.responseWithCustomMeta(c, "District retrieved successfully", responseData, page, limit, totalVillages)
	}
	return h.responseWithCustomMetaNoPage(c, "District retrieved successfully", responseData, totalVillages)
}

func (h *WilayahIndonesiaHandler) GetAllVillages(c *fiber.Ctx) error {
	ctx := c.Context()

	page, limit, isPaginated, err := h.parsePaginationParams(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	villages, total, err := h.WilayahService.GetAllVillages(ctx, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch villages: "+err.Error())
	}

	if isPaginated {
		return utils.SuccessWithPaginationAndTotal(c, "Villages retrieved successfully", villages, page, limit, total)
	}
	return utils.SuccessWithTotal(c, "Villages retrieved successfully", villages, total)
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

// parsePaginationParams returns page, limit, isPaginated, error
func (h *WilayahIndonesiaHandler) parsePaginationParams(c *fiber.Ctx) (int, int, bool, error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// If neither page nor limit is provided, return no pagination
	if pageStr == "" && limitStr == "" {
		return 0, 0, false, nil
	}

	// Default values when pagination is used
	page := 1
	limit := 10

	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, false, fiber.NewError(fiber.StatusBadRequest, "Invalid page parameter")
		}
		if parsedPage < 1 {
			return 0, 0, false, fiber.NewError(fiber.StatusBadRequest, "Page must be greater than 0")
		}
		page = parsedPage
	}

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, false, fiber.NewError(fiber.StatusBadRequest, "Invalid limit parameter")
		}
		if parsedLimit < 1 {
			return 0, 0, false, fiber.NewError(fiber.StatusBadRequest, "Limit must be greater than 0")
		}
		if parsedLimit > 100 {
			return 0, 0, false, fiber.NewError(fiber.StatusBadRequest, "Limit cannot exceed 100")
		}
		limit = parsedLimit
	}

	return page, limit, true, nil
}

// Custom response helpers for detailed responses with child counts
func (h *WilayahIndonesiaHandler) responseWithCustomMeta(c *fiber.Ctx, message string, data interface{}, page, limit, total int) error {
	type CustomMeta struct {
		Status int `json:"status"`
		Message string `json:"message"`
		Page   int `json:"page"`
		Limit  int `json:"limit"`
		Total  int `json:"total"`
	}

	type CustomResponse struct {
		Meta CustomMeta  `json:"meta"`
		Data interface{} `json:"data"`
	}

	response := CustomResponse{
		Meta: CustomMeta{
			Status:  200,
			Message: message,
			Page:    page,
			Limit:   limit,
			Total:   total,
		},
		Data: data,
	}

	return c.Status(200).JSON(response)
}

func (h *WilayahIndonesiaHandler) responseWithCustomMetaNoPage(c *fiber.Ctx, message string, data interface{}, total int) error {
	type CustomMeta struct {
		Status int `json:"status"`
		Message string `json:"message"`
		Total  int `json:"total"`
	}

	type CustomResponse struct {
		Meta CustomMeta  `json:"meta"`
		Data interface{} `json:"data"`
	}

	response := CustomResponse{
		Meta: CustomMeta{
			Status:  200,
			Message: message,
			Total:   total,
		},
		Data: data,
	}

	return c.Status(200).JSON(response)
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