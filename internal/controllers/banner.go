package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type BannerController struct {
	service *services.BannerService
}

func NewBannerController(service *services.BannerService) *BannerController {
	return &BannerController{service: service}
}

func (bc *BannerController) GetAllBanners(c *fiber.Ctx) error {
	banners, err := bc.service.GetAllBanners()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch banners",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banners fetched successfully",
		banners,
	))
}

func (bc *BannerController) GetBannerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	banner, err := bc.service.GetBannerByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
			fiber.StatusNotFound,
			"Banner not found",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banner fetched successfully",
		banner,
	))
}

func (bc *BannerController) CreateBanner(c *fiber.Ctx) error {
	var request dto.BannerCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	banner, err := bc.service.CreateBanner(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Banner created successfully",
		banner,
	))
}

func (bc *BannerController) UpdateBanner(c *fiber.Ctx) error {
	id := c.Params("id")
	var request dto.BannerUpdateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	banner, err := bc.service.UpdateBanner(id, &request)
	if err != nil {
		if err.Error() == "banner not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Banner not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banner updated successfully",
		banner,
	))
}

func (bc *BannerController) DeleteBanner(c *fiber.Ctx) error {
	id := c.Params("id")

	err := bc.service.DeleteBanner(id)
	if err != nil {
		if err.Error() == "banner not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Banner not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banner deleted successfully",
		nil,
	))
}
