package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetAllProducts(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 0)
	page := c.QueryInt("page", 1)

	if limit < 0 || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid pagination parameters",
			nil,
		))
	}

	var products []dto.ProductResponseDTO
	var err error

	if limit == 0 {
		products, err = services.GetProducts(0, 0)
	} else {
		products, err = services.GetProducts(limit, page)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch products",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Products fetched successfully",
		products,
	))
}

func GetProductByID(c *fiber.Ctx) error {
	productID := c.Params("productid")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Product ID is required",
			nil,
		))
	}

	product, err := services.GetProductByID(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Product not found",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Product fetched successfully",
		product,
	))
}
