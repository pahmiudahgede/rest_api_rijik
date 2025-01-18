package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetAllProducts(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: user ID is missing",
			nil,
		))
	}

	storeID := c.Query("storeID", "")

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	if limit <= 0 || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid pagination parameters",
			nil,
		))
	}

	var products []dto.ProductResponseWithSoldDTO
	var err error

	if storeID != "" {

		products, err = services.GetProductsByStoreID(storeID, limit, page)
	} else {

		products, err = services.GetProductsByUserID(userID, limit, page)
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
	storeID, ok := c.Locals("userID").(string)
	if !ok || storeID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: store ID is missing",
			nil,
		))
	}

	productID := c.Params("productid")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Product ID is required",
			nil,
		))
	}

	product, err := services.GetProductByIDAndStoreID(productID, storeID)
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

func GetProductsByStore(c *fiber.Ctx) error {
	storeID := c.Params("storeid")
	if storeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Store ID is required",
			nil,
		))
	}

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	if limit <= 0 || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid pagination parameters",
			nil,
		))
	}

	products, err := services.GetProductsByStoreID(storeID, limit, page)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Store not found",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Products fetched successfully",
		products,
	))
}

func CreateProduct(c *fiber.Ctx) error {
	var input dto.CreateProductRequestDTO

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request payload",
			nil,
		))
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: user ID is missing",
			nil,
		))
	}

	if input.StoreID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Store ID is required in the body",
			nil,
		))
	}

	product, err := services.CreateProduct(input, userID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.FormatResponse(
			fiber.StatusUnprocessableEntity,
			err.Error(),
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Product created successfully",
		product,
	))
}

func UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("productid")

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: user ID is missing",
			nil,
		))
	}

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Product ID is required",
			nil,
		))
	}

	var input dto.UpdateProductRequestDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request payload",
			nil,
		))
	}

	if err := dto.GetValidator().Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid product data",
			nil,
		))
	}

	updatedProduct, err := services.UpdateProduct(productID, input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.FormatResponse(
			fiber.StatusUnprocessableEntity,
			err.Error(),
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Product updated successfully",
		updatedProduct,
	))
}

func DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("productid")

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"Unauthorized: user ID is missing",
			nil,
		))
	}

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Product ID is required",
			nil,
		))
	}

	err := services.DeleteProduct(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Product not found or unable to delete",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Product deleted successfully",
		nil,
	))
}
