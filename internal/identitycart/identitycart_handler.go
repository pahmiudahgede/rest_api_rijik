package identitycart

import (
	"log"
	"rijig/middleware"
	"rijig/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type IdentityCardHandler struct {
	service IdentityCardService
}

func NewIdentityCardHandler(service IdentityCardService) *IdentityCardHandler {
	return &IdentityCardHandler{service: service}
}

func (h *IdentityCardHandler) CreateIdentityCardHandler(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "unauthorized access")
	}

	var input RequestIdentityCardDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.BadRequest(c, "Invalid input format")
	}

	if errs, valid := input.ValidateIdentityCardInput(); !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Input validation failed", errs)
	}

	cardPhoto, err := c.FormFile("cardphoto")
	if err != nil {
		log.Printf("Error getting card photo: %v", err)
		return utils.BadRequest(c, "KTP photo is required")
	}

	response, err := h.service.CreateIdentityCard(c.Context(), claims.UserID, claims.DeviceID, &input, cardPhoto)
	if err != nil {
		log.Printf("Error creating identity card: %v", err)
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create identity card")
	}

	return utils.SuccessWithData(c, "KTP successfully submitted", response)
}

func (h *IdentityCardHandler) GetIdentityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "ID is required")
	}

	result, err := h.service.GetIdentityCardByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting identity card by ID %s: %v", id, err)
		return utils.NotFound(c, "Identity card not found")
	}

	return utils.SuccessWithData(c, "Successfully retrieved identity card", result)
}

func (h *IdentityCardHandler) GetIdentityByUserId(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "Unauthorized access")
	}

	result, err := h.service.GetIdentityCardsByUserID(c.Context(), claims.UserID)
	if err != nil {
		log.Printf("Error getting identity cards for user %s: %v", claims.UserID, err)
		return utils.InternalServerError(c, "Failed to fetch your identity card data")
	}

	return utils.SuccessWithData(c, "Successfully retrieved your identity cards", result)
}

func (h *IdentityCardHandler) UpdateIdentityCardHandler(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "Unauthorized access")
	}

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Identity card ID is required")
	}

	var input RequestIdentityCardDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.BadRequest(c, "Invalid input format")
	}

	if errs, valid := input.ValidateIdentityCardInput(); !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Input validation failed", errs)
	}

	cardPhoto, err := c.FormFile("cardphoto")
	if err != nil && err.Error() != "there is no uploaded file associated with the given key" {
		log.Printf("Error getting card photo: %v", err)
		return utils.BadRequest(c, "Invalid card photo")
	}

	if cardPhoto != nil && cardPhoto.Size > 5*1024*1024 {
		return utils.BadRequest(c, "File size must be less than 5MB")
	}

	response, err := h.service.UpdateIdentityCard(c.Context(), claims.UserID, id, &input, cardPhoto)
	if err != nil {
		log.Printf("Error updating identity card: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Identity card not found")
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to update identity card")
	}

	return utils.SuccessWithData(c, "Identity card successfully updated", response)
}

func (h *IdentityCardHandler) GetAllIdentityCardsByRegStatus(c *fiber.Ctx) error {
	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "Unauthorized access")
	}

	// if claims.Role != "admin" {
	// 	return utils.Forbidden(c, "Access denied: admin role required")
	// }

	status := c.Query("status", utils.RegStatusPending)

	validStatuses := map[string]bool{
		utils.RegStatusPending: true,
		"confirmed":            true,
		"rejected":             true,
	}

	if !validStatuses[status] {
		return utils.BadRequest(c, "Invalid status. Valid values: pending, confirmed, rejected")
	}

	result, err := h.service.GetAllIdentityCardsByRegStatus(c.Context(), status)
	if err != nil {
		log.Printf("Error getting identity cards by status %s: %v", status, err)
		return utils.InternalServerError(c, "Failed to fetch identity cards")
	}

	return utils.SuccessWithData(c, "Successfully retrieved identity cards", result)
}

func (h *IdentityCardHandler) UpdateUserRegistrationStatusByIdentityCard(c *fiber.Ctx) error {
	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "Unauthorized access")
	}

	userID := c.Params("userId")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required")
	}

	type StatusUpdateRequest struct {
		Status string `json:"status" validate:"required,oneof=confirmed rejected"`
	}

	var input StatusUpdateRequest
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.BadRequest(c, "Invalid input format")
	}

	if input.Status != "confirmed" && input.Status != "rejected" {
		return utils.BadRequest(c, "Invalid status. Valid values: confirmed, rejected")
	}

	err = h.service.UpdateUserRegistrationStatusByIdentityCard(c.Context(), userID, input.Status)
	if err != nil {
		log.Printf("Error updating user registration status: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to update registration status")
	}

	message := "User registration status successfully updated to " + input.Status
	return utils.Success(c, message)
}

func (h *IdentityCardHandler) DeleteIdentityCardHandler(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Identity card ID is required")
	}

	return utils.Success(c, "Identity card successfully deleted")
}
