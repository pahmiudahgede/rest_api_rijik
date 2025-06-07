package identitycart

import (
	"rijig/middleware"
	"rijig/utils"

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
		return err
	}

	cardPhoto, err := c.FormFile("cardphoto")
	if err != nil {
		return utils.BadRequest(c, "KTP photo is required")
	}

	var input RequestIdentityCardDTO
	if err := c.BodyParser(&input); err != nil {
		return utils.BadRequest(c, "Invalid input format")
	}


	if errs, valid := input.ValidateIdentityCardInput(); !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Input validation failed", errs)
	}

	response, err := h.service.CreateIdentityCard(c.Context(), claims.UserID, &input, cardPhoto)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "KTP successfully submitted", response)
}

func (h *IdentityCardHandler) GetIdentityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "id is required")
	}

	result, err := h.service.GetIdentityCardByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "data not found")
	}

	return utils.SuccessWithData(c, "success retrieve identity card", result)
}

func (h *IdentityCardHandler) GetIdentityByUserId(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetIdentityCardsByUserID(c.Context(), claims.UserID)
	if err != nil {
		return utils.InternalServerError(c, "failed to fetch your identity card data")
	}

	return utils.SuccessWithData(c, "success retrieve your identity card", result)
}
