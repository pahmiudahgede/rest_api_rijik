package handler

import (
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

func (h *WilayahIndonesiaHandler) GetProvinces(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	provinces, totalProvinces, err := h.WilayahService.GetAllProvinces(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch provinces")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, provinces, page, limit, totalProvinces, "Provinces fetched successfully")
	}

	return utils.NonPaginatedResponse(c, provinces, totalProvinces, "Provinces fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetProvinceByID(c *fiber.Ctx) error {
	provinceID := c.Params("id")

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	province, totalRegencies, err := h.WilayahService.GetProvinceByID(provinceID, page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch province")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, province, page, limit, totalRegencies, "Province fetched successfully")
	}

	return utils.NonPaginatedResponse(c, province, totalRegencies, "Province fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllRegencies(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	regencies, totalRegencies, err := h.WilayahService.GetAllRegencies(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch regency")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, regencies, page, limit, totalRegencies, "regency fetched successfully")
	}

	return utils.NonPaginatedResponse(c, regencies, totalRegencies, "Provinces fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetRegencyByID(c *fiber.Ctx) error {
	regencyId := c.Params("id")

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	regency, totalDistrict, err := h.WilayahService.GetRegencyByID(regencyId, page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch regency")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, regency, page, limit, totalDistrict, "regency fetched successfully")
	}

	return utils.NonPaginatedResponse(c, regency, totalDistrict, "regency fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllDistricts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	districts, totalDistricts, err := h.WilayahService.GetAllDistricts(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch districts")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, districts, page, limit, totalDistricts, "districts fetched successfully")
	}

	return utils.NonPaginatedResponse(c, districts, totalDistricts, "districts fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetDistrictByID(c *fiber.Ctx) error {
	districtId := c.Params("id")

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	district, totalVillages, err := h.WilayahService.GetDistrictByID(districtId, page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch district")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, district, page, limit, totalVillages, "district fetched successfully")
	}

	return utils.NonPaginatedResponse(c, district, totalVillages, "district fetched successfully")
}

func (h *WilayahIndonesiaHandler) GetAllVillages(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	villages, totalVillages, err := h.WilayahService.GetAllVillages(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch villages")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, villages, page, limit, totalVillages, "villages fetched successfully")
	}

	return utils.NonPaginatedResponse(c, villages, totalVillages, "villages fetched successfully")
}
