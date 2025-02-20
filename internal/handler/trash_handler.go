package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type TrashHandler struct {
	TrashService services.TrashService
}

func NewTrashHandler(trashService services.TrashService) *TrashHandler {
	return &TrashHandler{TrashService: trashService}
}

func (h *TrashHandler) CreateCategory(c *fiber.Ctx) error {
	var request dto.RequestTrashCategoryDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	categoryResponse, err := h.TrashService.CreateCategory(request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to create category: "+err.Error())
	}

	return utils.CreateResponse(c, categoryResponse, "Category created successfully")
}

func (h *TrashHandler) AddDetailToCategory(c *fiber.Ctx) error {
	var request dto.RequestTrashDetailDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	detailResponse, err := h.TrashService.AddDetailToCategory(request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to add detail to category: "+err.Error())
	}

	return utils.CreateResponse(c, detailResponse, "Trash detail added successfully")
}

func (h *TrashHandler) GetCategories(c *fiber.Ctx) error {

	categories, err := h.TrashService.GetCategories()
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to fetch categories: "+err.Error())
	}

	return utils.NonPaginatedResponse(c, categories, len(categories), "Categories retrieved successfully")
}

func (h *TrashHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("category_id")

	category, err := h.TrashService.GetCategoryByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "Category not found: "+err.Error())
	}

	return utils.SuccessResponse(c, category, "Category retrieved successfully")
}

func (h *TrashHandler) GetTrashDetailByID(c *fiber.Ctx) error {
	id := c.Params("detail_id")

	detail, err := h.TrashService.GetTrashDetailByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "Trash detail not found: "+err.Error())
	}

	return utils.SuccessResponse(c, detail, "Trash detail retrieved successfully")
}

func (h *TrashHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("category_id")

	var request dto.RequestTrashCategoryDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid request body"}})
	}

	updatedCategory, err := h.TrashService.UpdateCategory(id, request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Error updating category: "+err.Error())
	}

	return utils.SuccessResponse(c, updatedCategory, "Category updated successfully")
}

func (h *TrashHandler) UpdateDetail(c *fiber.Ctx) error {
	id := c.Params("detail_id")

	var request dto.RequestTrashDetailDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid request body"}})
	}

	updatedDetail, err := h.TrashService.UpdateDetail(id, request)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Error updating detail: "+err.Error())
	}

	return utils.SuccessResponse(c, updatedDetail, "Trash detail updated successfully")
}

func (h *TrashHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("category_id")

	if err := h.TrashService.DeleteCategory(id); err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Error deleting category: "+err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Category deleted successfully")
}

func (h *TrashHandler) DeleteDetail(c *fiber.Ctx) error {
	id := c.Params("detail_id")

	if err := h.TrashService.DeleteDetail(id); err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Error deleting detail: "+err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Trash detail deleted successfully")
}
