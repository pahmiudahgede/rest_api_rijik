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

	errors, valid := requestUserPinDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	_, err := h.UserPinService.VerifyUserPin(requestUserPinDTO.Pin, userID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "pin yang anda masukkan salah")
	}

	return utils.LogResponse(c, map[string]string{"data": "pin yang anda masukkan benar"}, "Pin verification successful")
}

func (h *UserPinHandler) CheckPinStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	status, _, err := h.UserPinService.CheckPinStatus(userID)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if status == "Pin not created" {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, "pin belum dibuat")
	}

	return utils.LogResponse(c, map[string]string{"data": "pin sudah dibuat"}, "Pin status retrieved successfully")
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

	userPinResponse, err := h.UserPinService.CreateUserPin(userID, requestUserPinDTO.Pin)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusConflict, err.Error())
	}

	return utils.LogResponse(c, userPinResponse, "User pin created successfully")
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

	userPinResponse, err := h.UserPinService.UpdateUserPin(userID, requestUserPinDTO.OldPin, requestUserPinDTO.NewPin)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.LogResponse(c, userPinResponse, "User pin updated successfully")
}
