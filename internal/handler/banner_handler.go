package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type BannerHandler struct {
	BannerService services.BannerService
}

func NewBannerHandler(bannerService services.BannerService) *BannerHandler {
	return &BannerHandler{BannerService: bannerService}
}

func (h *BannerHandler) CreateBanner(c *fiber.Ctx) error {
	var request dto.RequestBannerDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateBannerInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	bannerImage, err := c.FormFile("bannerimage")
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Banner image is required")
	}

	bannerResponse, err := h.BannerService.CreateBanner(request, bannerImage)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.CreateResponse(c, bannerResponse, "Banner created successfully")
}

func (h *BannerHandler) GetAllBanners(c *fiber.Ctx) error {
	banners, err := h.BannerService.GetAllBanners()
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to fetch banners")
	}

	return utils.NonPaginatedResponse(c, banners, len(banners), "Banners fetched successfully")
}

func (h *BannerHandler) GetBannerByID(c *fiber.Ctx) error {
	id := c.Params("banner_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Banner ID is required")
	}

	banner, err := h.BannerService.GetBannerByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "invalid banner id")
	}

	return utils.SuccessResponse(c, banner, "Banner fetched successfully")
}

func (h *BannerHandler) UpdateBanner(c *fiber.Ctx) error {
	id := c.Params("banner_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Banner ID is required")
	}

	var request dto.RequestBannerDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateBannerInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	bannerImage, err := c.FormFile("bannerimage")
	if err != nil && err.Error() != "no such file" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Banner image is required")
	}

	bannerResponse, err := h.BannerService.UpdateBanner(id, request, bannerImage)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, bannerResponse, "Banner updated successfully")
}

func (h *BannerHandler) DeleteBanner(c *fiber.Ctx) error {
	id := c.Params("banner_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Banner ID is required")
	}

	err := h.BannerService.DeleteBanner(id)
	if err != nil {

		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Banner deleted successfully")
}
