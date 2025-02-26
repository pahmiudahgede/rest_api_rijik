package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
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
