package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetBanners(c *fiber.Ctx) error {

	banners, err := services.GetBanners()
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch banners",
			nil,
		))
	}

	var bannerResponses []dto.BannerResponse
	for _, banner := range banners {

		bannerResponses = append(bannerResponses, dto.BannerResponse{
			ID:          banner.ID,
			BannerName:  banner.BannerName,
			BannerImage: banner.BannerImage,
			CreatedAt:   utils.FormatDateToIndonesianFormat(banner.CreatedAt),
			UpdatedAt:   utils.FormatDateToIndonesianFormat(banner.UpdatedAt),
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banners fetched successfully",
		struct {
			Banners []dto.BannerResponse `json:"banners"`
		}{
			Banners: bannerResponses,
		},
	))
}

func GetBannerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	banner, err := services.GetBannerByID(id)
	if err != nil {

		if err.Error() == "banner not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
				fiber.StatusNotFound,
				"Banner not found",
				nil,
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch banner",
			nil,
		))
	}

	bannerResponse := dto.BannerResponse{
		ID:          banner.ID,
		BannerName:  banner.BannerName,
		BannerImage: banner.BannerImage,
		CreatedAt:   utils.FormatDateToIndonesianFormat(banner.CreatedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(banner.UpdatedAt),
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Banner fetched successfully",
		struct {
			Banner dto.BannerResponse `json:"banner"`
		}{
			Banner: bannerResponse,
		},
	))
}

func CreateBanner(c *fiber.Ctx) error {
	var bannerInput dto.BannerRequest

	if err := c.BodyParser(&bannerInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := bannerInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	newBanner, err := services.CreateBanner(bannerInput.BannerName, bannerInput.BannerImage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create banner",
			nil,
		))
	}

	bannerResponse := dto.NewBannerResponse(
		newBanner.ID,
		newBanner.BannerName,
		newBanner.BannerImage,
		utils.FormatDateToIndonesianFormat(newBanner.CreatedAt),
		utils.FormatDateToIndonesianFormat(newBanner.UpdatedAt),
	)

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Banner created successfully",
		struct {
			Banner dto.BannerResponse `json:"banner"`
		}{
			Banner: bannerResponse,
		},
	))
}