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

	var requestStoreDTO dto.RequestStoreDTO
	if err := c.BodyParser(&requestStoreDTO); err != nil {

		log.Printf("Error parsing body: %v", err)
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid request body"}})
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

	store, err := h.StoreService.CreateStore(userID, &requestStoreDTO)
	if err != nil {

		log.Printf("Error creating store: %v", err)
		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.CreateResponse(c, store, "store created successfully")
}
