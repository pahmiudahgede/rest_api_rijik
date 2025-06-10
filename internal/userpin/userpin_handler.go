package userpin

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type UserPinHandler struct {
	service UserPinService
}

func NewUserPinHandler(service UserPinService) *UserPinHandler {
	return &UserPinHandler{service}
}

func (h *UserPinHandler) CreateUserPinHandler(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Authentication required")
	}

	if claims.UserID == "" || claims.DeviceID == "" {
		return utils.BadRequest(c, "Invalid user claims")
	}

	var req RequestPinDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if errs, ok := req.ValidateRequestPinDTO(); !ok {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation error", errs)
	}

	pintokenresponse, err := h.service.CreateUserPin(c.Context(), claims.UserID, claims.DeviceID, &req)
	if err != nil {
		if err.Error() == Pinhasbeencreated {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "PIN created successfully", pintokenresponse)
}

func (h *UserPinHandler) VerifyPinHandler(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req RequestPinDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	token, err := h.service.VerifyUserPin(c.Context(), claims.UserID, claims.DeviceID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "PIN verified successfully", token)
}
