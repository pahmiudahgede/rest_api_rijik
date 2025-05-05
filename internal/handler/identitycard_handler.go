package handler

import (
	"log"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type IdentityCardHandler struct {
	IdentityCardService services.IdentityCardService
}

func NewIdentityCardHandler(identityCardService services.IdentityCardService) *IdentityCardHandler {
	return &IdentityCardHandler{
		IdentityCardService: identityCardService,
	}
}

func (h *IdentityCardHandler) CreateIdentityCard(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	var request dto.RequestIdentityCardDTO
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.ErrorResponse(c, "Invalid request data")
	}

	cardPhoto, err := c.FormFile("cardphoto")
	if err != nil {
		log.Printf("Error retrieving card photo from request: %v", err)
		return utils.ErrorResponse(c, "Card photo is required")
	}

	identityCard, err := h.IdentityCardService.CreateIdentityCard(userID, &request, cardPhoto)
	if err != nil {
		log.Printf("Error creating identity card: %v", err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.CreateResponse(c, identityCard, "Identity card created successfully")
}

func (h *IdentityCardHandler) UpdateIdentityCard(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	id := c.Params("identity_id")
	if id == "" {
		return utils.ErrorResponse(c, "Identity card ID is required")
	}

	var request dto.RequestIdentityCardDTO
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.ErrorResponse(c, "Invalid request data")
	}

	cardPhoto, err := c.FormFile("cardphoto")
	if err != nil && err.Error() != "File not found" {
		log.Printf("Error retrieving card photo: %v", err)
		return utils.ErrorResponse(c, "Card photo is required")
	}

	updatedCard, err := h.IdentityCardService.UpdateIdentityCard(userID, id, &request, cardPhoto)
	if err != nil {
		log.Printf("Error updating identity card: %v", err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, updatedCard, "Identity card updated successfully")
}

func (h *IdentityCardHandler) GetIdentityCardById(c *fiber.Ctx) error {

	id := c.Params("identity_id")
	if id == "" {
		return utils.ErrorResponse(c, "Identity card ID is required")
	}

	identityCard, err := h.IdentityCardService.GetIdentityCardByID(id)
	if err != nil {
		log.Printf("Error retrieving identity card: %v", err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, identityCard, "Identity card retrieved successfully")
}

func (h *IdentityCardHandler) GetIdentityCard(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	identityCard, err := h.IdentityCardService.GetIdentityCardsByUserID(userID)
	if err != nil {
		log.Printf("Error retrieving identity card: %v", err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, identityCard, "Identity card retrieved successfully")
}

func (h *IdentityCardHandler) DeleteIdentityCard(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	id := c.Params("identity_id")
	if id == "" {
		return utils.ErrorResponse(c, "Identity card ID is required")
	}

	err := h.IdentityCardService.DeleteIdentityCard(id)
	if err != nil {
		log.Printf("Error deleting identity card: %v", err)
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Identity card deleted successfully")
}
