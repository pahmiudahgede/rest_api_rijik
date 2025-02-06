package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type WilayahIndonesiaHandler struct {
	WilayahService services.WilayahIndonesiaService
}

func NewWilayahImportHandler(wilayahService services.WilayahIndonesiaService) *WilayahIndonesiaHandler {
	return &WilayahIndonesiaHandler{WilayahService: wilayahService}
}

func (h *WilayahIndonesiaHandler) ImportWilayahData(c *fiber.Ctx) error {

	err := h.WilayahService.ImportDataFromCSV()
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.GenericErrorResponse(c, fiber.StatusCreated, "Data imported successfully")
}

func (h *WilayahIndonesiaHandler) GetAllProvinces(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	provinces, err := h.WilayahService.GetAllProvinces(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch provinces: %v", err))
	}

	return utils.PaginatedResponse(c, provinces, page, limit, len(provinces), "Provinces fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetProvinceByID(c *fiber.Ctx) error {
	provinceID := c.Params("id")

	province, err := h.WilayahService.GetProvinceByID(provinceID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("Province not found: %v", err))
	}

	return utils.LogResponse(c, province, "Province fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllRegencies(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	regencies, err := h.WilayahService.GetAllRegencies(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch regencies: %v", err))
	}

	return utils.PaginatedResponse(c, regencies, page, limit, len(regencies), "Regencies fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetRegencyByID(c *fiber.Ctx) error {
	regencyID := c.Params("id")

	regency, err := h.WilayahService.GetRegencyByID(regencyID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("Regency not found: %v", err))
	}

	return utils.LogResponse(c, regency, "Regency fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllDistricts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	districts, err := h.WilayahService.GetAllDistricts(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch districts: %v", err))
	}

	return utils.PaginatedResponse(c, districts, page, limit, len(districts), "Districts fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetDistrictByID(c *fiber.Ctx) error {
	districtID := c.Params("id")

	district, err := h.WilayahService.GetDistrictByID(districtID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("District not found: %v", err))
	}

	return utils.LogResponse(c, district, "District fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllVillages(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	villages, err := h.WilayahService.GetAllVillages(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch villages: %v", err))
	}

	return utils.PaginatedResponse(c, villages, page, limit, len(villages), "Villages fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetVillageByID(c *fiber.Ctx) error {
	villageID := c.Params("id")

	village, err := h.WilayahService.GetVillageByID(villageID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("Village not found: %v", err))
	}

	return utils.LogResponse(c, village, "Village fetched successfully")
}
