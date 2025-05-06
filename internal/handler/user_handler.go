package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UpdateUserAvatarHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	avatar, err := c.FormFile("avatar")
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "No avatar file provided")
	}

	updatedUser, err := h.userService.UpdateUserAvatar(userID, avatar)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, updatedUser, "Avatar updated successfully")
}

func (h *UserHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	// userID := c.Params("id")
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, user, "User retrieved successfully")
}

func (h *UserHandler) GetAllUsersHandler(c *fiber.Ctx) error {

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	users, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.PaginatedResponse(c, users, page, limit, len(users), "Users retrieved successfully")
}

func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	var request dto.RequestUserDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userID").(string)
	updatedUser, err := h.userService.UpdateUser(userID, &request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, updatedUser, "User profile updated successfully")
}

func (h *UserHandler) UpdateUserPasswordHandler(c *fiber.Ctx) error {
	var request dto.UpdatePasswordDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userID").(string)
	err := h.userService.UpdateUserPassword(userID, request.OldPassword, request.NewPassword, request.ConfirmNewPassword)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Password updated successfully")
}
