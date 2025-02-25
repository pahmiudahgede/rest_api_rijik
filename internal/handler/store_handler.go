package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type StoreHandler struct {
	StoreService services.StoreService
}

func NewStoreHandler(storeService services.StoreService) *StoreHandler {
	return &StoreHandler{StoreService: storeService}
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {

	storeName := c.FormValue("store_name")
	storeInfo := c.FormValue("store_info")
	storeAddressID := c.FormValue("store_address_id")

	if storeName == "" || storeInfo == "" || storeAddressID == "" {
		log.Println("Missing required fields")
		return utils.GenericResponse(c, fiber.StatusBadRequest, "All fields are required")
	}

	storeLogo, err := c.FormFile("store_logo")
	if err != nil {
		log.Printf("Error parsing store logo: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Store logo is required")
	}

	storeBanner, err := c.FormFile("store_banner")
	if err != nil {
		log.Printf("Error parsing store banner: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Store banner is required")
	}

	requestStoreDTO := dto.RequestStoreDTO{
		StoreName:      storeName,
		StoreLogo:      storeLogo.Filename,
		StoreBanner:    storeBanner.Filename,
		StoreInfo:      storeInfo,
		StoreAddressID: storeAddressID,
	}

	errors, valid := requestStoreDTO.ValidateStoreInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	store, err := h.StoreService.CreateStore(userID, requestStoreDTO, storeLogo, storeBanner)
	if err != nil {
		log.Printf("Error creating store: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.CreateResponse(c, store, "Store created successfully")
}

func (h *StoreHandler) GetStoreByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	store, err := h.StoreService.GetStoreByUserID(userID)
	if err != nil {
		log.Printf("Error fetching store: %v", err)
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	log.Printf("Store fetched successfully: %v", store)
	return utils.SuccessResponse(c, store, "Store fetched successfully")
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	storeID := c.Params("store_id")
	if storeID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Store ID is required")
	}

	storeName := c.FormValue("store_name")
	storeInfo := c.FormValue("store_info")
	storeAddressID := c.FormValue("store_address_id")

	if storeName == "" || storeInfo == "" || storeAddressID == "" {
		log.Println("Missing required fields")
		return utils.GenericResponse(c, fiber.StatusBadRequest, "All fields are required")
	}

	storeLogo, err := c.FormFile("store_logo")
	if err != nil && err.Error() != "missing file" {
		log.Printf("Error parsing store logo: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Error parsing store logo")
	}

	storeBanner, err := c.FormFile("store_banner")
	if err != nil && err.Error() != "missing file" {
		log.Printf("Error parsing store banner: %v", err)
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Error parsing store banner")
	}

	requestStoreDTO := dto.RequestStoreDTO{
		StoreName:      storeName,
		StoreLogo:      storeLogo.Filename,
		StoreBanner:    storeBanner.Filename,
		StoreInfo:      storeInfo,
		StoreAddressID: storeAddressID,
	}

	errors, valid := requestStoreDTO.ValidateStoreInput()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("User ID not found in Locals")
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User ID not found")
	}

	store, err := h.StoreService.UpdateStore(storeID, &requestStoreDTO, storeLogo, storeBanner, userID)
	if err != nil {
		log.Printf("Error updating store: %v", err)
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	log.Printf("Store updated successfully: %v", store)
	return utils.SuccessResponse(c, store, "Store updated successfully")
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	storeID := c.Params("store_id")
	if storeID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Store ID is required")
	}

	err := h.StoreService.DeleteStore(storeID)
	if err != nil {
		log.Printf("Error deleting store: %v", err)
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	log.Printf("Store deleted successfully: %v", storeID)
	return utils.GenericResponse(c, fiber.StatusOK, "Store deleted successfully")
}
