package handler

import (
	"fmt"
	"log"
	"strconv"

	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func ConvertStringToInt(value string) (int, error) {
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer format: %s", value)
	}
	return convertedValue, nil
}

func GetPaginationParams(c *fiber.Ctx) (int, int, error) {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 0, 0, fmt.Errorf("invalid page value")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, fmt.Errorf("invalid limit value")
	}

	return page, limit, nil
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	productName := c.FormValue("product_name")
	quantityStr := c.FormValue("quantity")
	productImages, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Error parsing form data")
	}

	quantity, err := ConvertStringToInt(quantityStr)
	if err != nil {
		log.Printf("Invalid quantity: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid quantity")
	}

	productDTO := dto.RequestProductDTO{
		ProductName:   productName,
		Quantity:      quantity,
		ProductImages: productImages.File["product_image"],
	}

	product, err := h.ProductService.CreateProduct(userID, &productDTO)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.CreateResponse(c, product, "Product created successfully")
}

func (h *ProductHandler) GetAllProductsByStoreID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	page, limit, err := GetPaginationParams(c)
	if err != nil {
		log.Printf("Invalid pagination params: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters")
	}

	products, total, err := h.ProductService.GetAllProductsByStoreID(userID, page, limit)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.PaginatedResponse(c, products, page, limit, int(total), "Products fetched successfully")
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {

	productID := c.Params("product_id")
	if productID == "" {
		log.Println("Product ID is required")
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Product ID is required")
	}

	product, err := h.ProductService.GetProductByID(productID)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, product, "Product fetched successfully")
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	productID := c.Params("product_id")
	if productID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Product ID is required")
	}

	var productDTO dto.RequestProductDTO
	if err := c.BodyParser(&productDTO); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	productImages, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Error parsing form data")
	}

	productDTO.ProductImages = productImages.File["product_images"]

	product, err := h.ProductService.UpdateProduct(userID, productID, &productDTO)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.CreateResponse(c, product, "Product updated successfully")
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if productID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Product ID is required")
	}

	err := h.ProductService.DeleteProduct(productID)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, fmt.Sprintf("Failed to delete product: %v", err))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Product deleted successfully")
}

func (h *ProductHandler) DeleteProducts(c *fiber.Ctx) error {
	var productIDs []string
	if err := c.BodyParser(&productIDs); err != nil {
		log.Printf("Error parsing product IDs: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid input for product IDs")
	}

	if len(productIDs) == 0 {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "No product IDs provided")
	}

	err := h.ProductService.DeleteProducts(productIDs)
	if err != nil {
		log.Printf("Error deleting products: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, fmt.Sprintf("Failed to delete products: %v", err))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Products deleted successfully")
}

func (h *ProductHandler) DeleteProductImage(c *fiber.Ctx) error {
	imageID := c.Params("image_id")
	if imageID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Image ID is required")
	}

	err := h.ProductService.DeleteProductImage(imageID)
	if err != nil {
		log.Printf("Error deleting product image: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, fmt.Sprintf("Failed to delete product image: %v", err))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Product image deleted successfully")
}

func (h *ProductHandler) DeleteProductImages(c *fiber.Ctx) error {
	var imageIDs []string
	if err := c.BodyParser(&imageIDs); err != nil {
		log.Printf("Error parsing image IDs: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid input for image IDs")
	}

	if len(imageIDs) == 0 {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "No image IDs provided")
	}

	err := h.ProductService.DeleteProductImages(imageIDs)
	if err != nil {
		log.Printf("Error deleting product images: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, fmt.Sprintf("Failed to delete product images: %v", err))
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Product images deleted successfully")
}
