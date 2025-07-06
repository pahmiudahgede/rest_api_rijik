package trash

import (
	"rijig/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TrashHandler struct {
	trashService TrashServiceInterface
}

func NewTrashHandler(trashService TrashServiceInterface) *TrashHandler {
	return &TrashHandler{
		trashService: trashService,
	}
}

func (h *TrashHandler) CreateTrashCategory(c *fiber.Ctx) error {
	var req RequestTrashCategoryDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.CreateTrashCategory(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		return utils.InternalServerError(c, "Failed to create trash category")
	}

	return utils.CreateSuccessWithData(c, "Trash category created successfully", response)
}

func (h *TrashHandler) CreateTrashCategoryWithIcon(c *fiber.Ctx) error {
	var req RequestTrashCategoryDTO

	req.Name = c.FormValue("name")
	req.Variety = c.FormValue("variety")

	if estimatedPriceStr := c.FormValue("estimated_price"); estimatedPriceStr != "" {
		if price, err := strconv.ParseFloat(estimatedPriceStr, 64); err == nil {
			req.EstimatedPrice = price
		}
	}

	iconFile, err := c.FormFile("icon")
	if err != nil && err.Error() != "there is no uploaded file associated with the given key" {
		return utils.BadRequest(c, "Invalid icon file")
	}

	response, err := h.trashService.CreateTrashCategoryWithIcon(c.Context(), req, iconFile)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create trash category")
	}

	return utils.CreateSuccessWithData(c, "Trash category created successfully", response)
}

func (h *TrashHandler) CreateTrashCategoryWithDetails(c *fiber.Ctx) error {
	var req struct {
		Category RequestTrashCategoryDTO `json:"category"`
		Details  []RequestTrashDetailDTO `json:"details"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.CreateTrashCategoryWithDetails(c.Context(), req.Category, req.Details)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		return utils.InternalServerError(c, "Failed to create trash category with details")
	}

	return utils.CreateSuccessWithData(c, "Trash category with details created successfully", response)
}

func (h *TrashHandler) UpdateTrashCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req RequestTrashCategoryDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.UpdateTrashCategory(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to update trash category")
	}

	return utils.SuccessWithData(c, "Trash category updated successfully", response)
}

func (h *TrashHandler) UpdateTrashCategoryWithIcon(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req RequestTrashCategoryDTO

	req.Name = c.FormValue("name")
	req.Variety = c.FormValue("variety")

	if estimatedPriceStr := c.FormValue("estimated_price"); estimatedPriceStr != "" {
		if price, err := strconv.ParseFloat(estimatedPriceStr, 64); err == nil {
			req.EstimatedPrice = price
		}
	}

	iconFile, _ := c.FormFile("icon")

	response, err := h.trashService.UpdateTrashCategoryWithIcon(c.Context(), id, req, iconFile)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to update trash category")
	}

	return utils.SuccessWithData(c, "Trash category updated successfully", response)
}

func (h *TrashHandler) GetAllTrashCategories(c *fiber.Ctx) error {
	withDetails := c.Query("with_details", "false")

	if withDetails == "true" {
		response, err := h.trashService.GetAllTrashCategoriesWithDetails(c.Context())
		if err != nil {
			return utils.InternalServerError(c, "Failed to get trash categories")
		}
		return utils.SuccessWithData(c, "Trash categories retrieved successfully", response)
	}

	response, err := h.trashService.GetAllTrashCategories(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to get trash categories")
	}

	return utils.SuccessWithData(c, "Trash categories retrieved successfully", response)
}

func (h *TrashHandler) GetTrashCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	withDetails := c.Query("with_details", "false")

	if withDetails == "true" {
		response, err := h.trashService.GetTrashCategoryByIDWithDetails(c.Context(), id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return utils.NotFound(c, "Trash category not found")
			}
			return utils.InternalServerError(c, "Failed to get trash category")
		}
		return utils.SuccessWithData(c, "Trash category retrieved successfully", response)
	}

	response, err := h.trashService.GetTrashCategoryByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to get trash category")
	}

	return utils.SuccessWithData(c, "Trash category retrieved successfully", response)
}

func (h *TrashHandler) DeleteTrashCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	err := h.trashService.DeleteTrashCategory(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to delete trash category")
	}

	return utils.Success(c, "Trash category deleted successfully")
}

func (h *TrashHandler) CreateTrashDetail(c *fiber.Ctx) error {
	var req RequestTrashDetailDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.CreateTrashDetail(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		return utils.InternalServerError(c, "Failed to create trash detail")
	}

	return utils.CreateSuccessWithData(c, "Trash detail created successfully", response)
}

func (h *TrashHandler) CreateTrashDetailWithIcon(c *fiber.Ctx) error {
	var req RequestTrashDetailDTO

	req.CategoryID = c.FormValue("category_id")
	req.Description = c.FormValue("description")

	if stepOrderStr := c.FormValue("step_order"); stepOrderStr != "" {
		if stepOrder, err := strconv.Atoi(stepOrderStr); err == nil {
			req.StepOrder = stepOrder
		}
	}

	iconFile, err := c.FormFile("icon")
	if err != nil && err.Error() != "there is no uploaded file associated with the given key" {
		return utils.BadRequest(c, "Invalid icon file")
	}

	response, err := h.trashService.CreateTrashDetailWithIcon(c.Context(), req, iconFile)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create trash detail")
	}

	return utils.CreateSuccessWithData(c, "Trash detail created successfully", response)
}

func (h *TrashHandler) AddTrashDetailToCategory(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if categoryID == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req RequestTrashDetailDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.AddTrashDetailToCategory(c.Context(), categoryID, req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		return utils.InternalServerError(c, "Failed to add trash detail to category")
	}

	return utils.CreateSuccessWithData(c, "Trash detail added to category successfully", response)
}

func (h *TrashHandler) AddTrashDetailToCategoryWithIcon(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if categoryID == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req RequestTrashDetailDTO

	req.Description = c.FormValue("description")

	if stepOrderStr := c.FormValue("step_order"); stepOrderStr != "" {
		if stepOrder, err := strconv.Atoi(stepOrderStr); err == nil {
			req.StepOrder = stepOrder
		}
	}

	iconFile, err := c.FormFile("icon")
	if err != nil && err.Error() != "there is no uploaded file associated with the given key" {
		return utils.BadRequest(c, "Invalid icon file")
	}

	response, err := h.trashService.AddTrashDetailToCategoryWithIcon(c.Context(), categoryID, req, iconFile)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to add trash detail to category")
	}

	return utils.CreateSuccessWithData(c, "Trash detail added to category successfully", response)
}

func (h *TrashHandler) UpdateTrashDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Detail ID is required")
	}

	var req RequestTrashDetailDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	response, err := h.trashService.UpdateTrashDetail(c.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash detail not found")
		}
		return utils.InternalServerError(c, "Failed to update trash detail")
	}

	return utils.SuccessWithData(c, "Trash detail updated successfully", response)
}

func (h *TrashHandler) UpdateTrashDetailWithIcon(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Detail ID is required")
	}

	var req RequestTrashDetailDTO

	req.Description = c.FormValue("description")

	if stepOrderStr := c.FormValue("step_order"); stepOrderStr != "" {
		if stepOrder, err := strconv.Atoi(stepOrderStr); err == nil {
			req.StepOrder = stepOrder
		}
	}

	iconFile, _ := c.FormFile("icon")

	response, err := h.trashService.UpdateTrashDetailWithIcon(c.Context(), id, req, iconFile)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash detail not found")
		}
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to update trash detail")
	}

	return utils.SuccessWithData(c, "Trash detail updated successfully", response)
}

func (h *TrashHandler) GetTrashDetailsByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if categoryID == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	response, err := h.trashService.GetTrashDetailsByCategory(c.Context(), categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to get trash details")
	}

	return utils.SuccessWithData(c, "Trash details retrieved successfully", response)
}

func (h *TrashHandler) GetTrashDetailByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Detail ID is required")
	}

	response, err := h.trashService.GetTrashDetailByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash detail not found")
		}
		return utils.InternalServerError(c, "Failed to get trash detail")
	}

	return utils.SuccessWithData(c, "Trash detail retrieved successfully", response)
}

func (h *TrashHandler) DeleteTrashDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "Detail ID is required")
	}

	err := h.trashService.DeleteTrashDetail(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash detail not found")
		}
		return utils.InternalServerError(c, "Failed to delete trash detail")
	}

	return utils.Success(c, "Trash detail deleted successfully")
}

func (h *TrashHandler) BulkCreateTrashDetails(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if categoryID == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req struct {
		Details []RequestTrashDetailDTO `json:"details"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if len(req.Details) == 0 {
		return utils.BadRequest(c, "At least one detail is required")
	}

	response, err := h.trashService.BulkCreateTrashDetails(c.Context(), categoryID, req.Details)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", extractValidationErrors(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to bulk create trash details")
	}

	return utils.CreateSuccessWithData(c, "Trash details created successfully", response)
}

func (h *TrashHandler) BulkDeleteTrashDetails(c *fiber.Ctx) error {
	var req struct {
		DetailIDs []string `json:"detail_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if len(req.DetailIDs) == 0 {
		return utils.BadRequest(c, "At least one detail ID is required")
	}

	err := h.trashService.BulkDeleteTrashDetails(c.Context(), req.DetailIDs)
	if err != nil {
		return utils.InternalServerError(c, "Failed to bulk delete trash details")
	}

	return utils.Success(c, "Trash details deleted successfully")
}

func (h *TrashHandler) ReorderTrashDetails(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if categoryID == "" {
		return utils.BadRequest(c, "Category ID is required")
	}

	var req struct {
		OrderedDetailIDs []string `json:"ordered_detail_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if len(req.OrderedDetailIDs) == 0 {
		return utils.BadRequest(c, "At least one detail ID is required")
	}

	err := h.trashService.ReorderTrashDetails(c.Context(), categoryID, req.OrderedDetailIDs)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Trash category not found")
		}
		return utils.InternalServerError(c, "Failed to reorder trash details")
	}

	return utils.Success(c, "Trash details reordered successfully")
}

func extractValidationErrors(errMsg string) interface{} {

	if strings.Contains(errMsg, "validation failed:") {
		return strings.TrimSpace(strings.Split(errMsg, "validation failed:")[1])
	}
	return errMsg
}
