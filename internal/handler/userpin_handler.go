package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type UserPinHandler struct {
	UserPinService services.UserPinService
}

func NewUserPinHandler(userPinService services.UserPinService) *UserPinHandler {
	return &UserPinHandler{UserPinService: userPinService}
}

func (h *UserPinHandler) VerifyUserPin(c *fiber.Ctx) error {
	var requestUserPinDTO dto.RequestUserPinDTO
	if err := c.BodyParser(&requestUserPinDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	// Validasi input pin
	errors, valid := requestUserPinDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	message, err := h.UserPinService.VerifyUserPin(userID, requestUserPinDTO.Pin)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, message)
}

func (h *UserPinHandler) CheckPinStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	status, err := h.UserPinService.CheckPinStatus(userID)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if status == "Pin not created" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Pin belum dibuat")
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Pin sudah dibuat")
}

func (h *UserPinHandler) CreateUserPin(c *fiber.Ctx) error {
	var requestUserPinDTO dto.RequestUserPinDTO
	if err := c.BodyParser(&requestUserPinDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := requestUserPinDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID := c.Locals("userID").(string)

	message, err := h.UserPinService.CreateUserPin(userID, requestUserPinDTO.Pin)
	if err != nil {

		return utils.GenericResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusCreated, message)
}

func (h *UserPinHandler) UpdateUserPin(c *fiber.Ctx) error {
	var requestUserPinDTO dto.UpdateUserPinDTO
	if err := c.BodyParser(&requestUserPinDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := requestUserPinDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID := c.Locals("userID").(string)

	message, err := h.UserPinService.UpdateUserPin(userID, requestUserPinDTO.OldPin, requestUserPinDTO.NewPin)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, message)
}
