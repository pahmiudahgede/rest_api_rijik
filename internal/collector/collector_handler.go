package collector

import (
	"rijig/middleware"
	"rijig/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CollectorHandler struct {
	collectorService CollectorService
}

func NewCollectorHandler(collectorService CollectorService) *CollectorHandler {
	return &CollectorHandler{
		collectorService: collectorService,
	}
}

func (h *CollectorHandler) CreateCollector(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var req CreateCollectorRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	errors, isValid := req.ValidateCreateCollectorRequest()
	if !isValid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	req.SetDefaults()

	collector, err := h.collectorService.CreateCollector(c.Context(), &req, claims.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create collector")
	}

	return utils.CreateSuccessWithData(c, "Collector created successfully", collector)
}

func (h *CollectorHandler) GetCollectorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return utils.BadRequest(c, "Collector ID is required")
	}

	collector, err := h.collectorService.GetCollectorByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to get collector")
	}

	return utils.SuccessWithData(c, "Collector retrieved successfully", collector)
}

func (h *CollectorHandler) GetCollectorByUserID(c *fiber.Ctx) error {
	// userID := c.Params("userID")
	// if strings.TrimSpace(userID) == "" {
	// 	return utils.BadRequest(c, "User ID is required")
	// }
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	collector, err := h.collectorService.GetCollectorByUserID(c.Context(), claims.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found for this user")
		}
		return utils.InternalServerError(c, "Failed to get collector")
	}

	return utils.SuccessWithData(c, "Collector retrieved successfully", collector)
}

func (h *CollectorHandler) UpdateCollector(c *fiber.Ctx) error {
	/* id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return utils.BadRequest(c, "Collector ID is required")
	} */
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req UpdateCollectorRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	errors, isValid := req.ValidateUpdateCollectorRequest()
	if !isValid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	req.NormalizeJobStatus()

	collector, err := h.collectorService.UpdateCollector(c.Context(), claims.UserID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to update collector")
	}

	return utils.SuccessWithData(c, "Collector updated successfully", collector)
}

func (h *CollectorHandler) DeleteCollector(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	err = h.collectorService.DeleteCollector(c.Context(), claims.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to delete collector")
	}

	return utils.Success(c, "Collector deleted successfully")
}

func (h *CollectorHandler) ListCollectors(c *fiber.Ctx) error {

	limit, offset, page := h.parsePaginationParams(c)

	collectors, total, err := h.collectorService.ListCollectors(c.Context(), limit, offset)
	if err != nil {
		return utils.InternalServerError(c, "Failed to get collectors")
	}

	responseData := map[string]interface{}{
		"collectors": collectors,
		"total":      total,
	}

	return utils.SuccessWithPagination(c, "Collectors retrieved successfully", responseData, page, limit)
}

func (h *CollectorHandler) GetActiveCollectors(c *fiber.Ctx) error {

	limit, offset, page := h.parsePaginationParams(c)

	collectors, total, err := h.collectorService.GetActiveCollectors(c.Context(), limit, offset)
	if err != nil {
		return utils.InternalServerError(c, "Failed to get active collectors")
	}

	responseData := map[string]interface{}{
		"collectors": collectors,
		"total":      total,
	}

	return utils.SuccessWithPagination(c, "Active collectors retrieved successfully", responseData, page, limit)
}

func (h *CollectorHandler) GetCollectorsByAddress(c *fiber.Ctx) error {
	addressID := c.Params("addressID")
	if strings.TrimSpace(addressID) == "" {
		return utils.BadRequest(c, "Address ID is required")
	}

	limit, offset, page := h.parsePaginationParams(c)

	collectors, total, err := h.collectorService.GetCollectorsByAddress(c.Context(), addressID, limit, offset)
	if err != nil {
		return utils.InternalServerError(c, "Failed to get collectors by address")
	}

	responseData := map[string]interface{}{
		"collectors": collectors,
		"total":      total,
		"address_id": addressID,
	}

	return utils.SuccessWithPagination(c, "Collectors by address retrieved successfully", responseData, page, limit)
}

func (h *CollectorHandler) GetCollectorsByTrashCategory(c *fiber.Ctx) error {
	trashCategoryID := c.Params("trashCategoryID")
	if strings.TrimSpace(trashCategoryID) == "" {
		return utils.BadRequest(c, "Trash category ID is required")
	}

	limit, offset, page := h.parsePaginationParams(c)

	collectors, total, err := h.collectorService.GetCollectorsByTrashCategory(c.Context(), trashCategoryID, limit, offset)
	if err != nil {
		return utils.InternalServerError(c, "Failed to get collectors by trash category")
	}

	responseData := map[string]interface{}{
		"collectors":        collectors,
		"total":             total,
		"trash_category_id": trashCategoryID,
	}

	return utils.SuccessWithPagination(c, "Collectors by trash category retrieved successfully", responseData, page, limit)
}

func (h *CollectorHandler) UpdateJobStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return utils.BadRequest(c, "Collector ID is required")
	}

	var req struct {
		JobStatus string `json:"job_status" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if strings.TrimSpace(req.JobStatus) == "" {
		return utils.BadRequest(c, "Job status is required")
	}

	jobStatus := strings.ToLower(strings.TrimSpace(req.JobStatus))
	validStatuses := []string{"active", "inactive", "busy"}
	if !h.isValidJobStatus(jobStatus, validStatuses) {
		return utils.BadRequest(c, "Invalid job status. Valid statuses: active, inactive, busy")
	}

	err := h.collectorService.UpdateJobStatus(c.Context(), id, jobStatus)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to update job status")
	}

	return utils.Success(c, "Job status updated successfully")
}

func (h *CollectorHandler) UpdateRating(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return utils.BadRequest(c, "Collector ID is required")
	}

	var req struct {
		Rating float32 `json:"rating" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if req.Rating < 1.0 || req.Rating > 5.0 {
		return utils.BadRequest(c, "Rating must be between 1.0 and 5.0")
	}

	err := h.collectorService.UpdateRating(c.Context(), id, req.Rating)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to update rating")
	}

	return utils.Success(c, "Rating updated successfully")
}

func (h *CollectorHandler) UpdateAvailableTrash(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return utils.BadRequest(c, "Collector ID is required")
	}

	var req BulkUpdateAvailableTrashRequest
	req.CollectorID = id

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	errors, isValid := req.ValidateBulkUpdateAvailableTrashRequest()
	if !isValid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	err := h.collectorService.UpdateAvailableTrash(c.Context(), id, req.AvailableTrashItems)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.NotFound(c, "Collector not found")
		}
		return utils.InternalServerError(c, "Failed to update available trash")
	}

	return utils.Success(c, "Available trash updated successfully")
}

func (h *CollectorHandler) parsePaginationParams(c *fiber.Ctx) (limit, offset, page int) {

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	pageStr := c.Query("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset = (page - 1) * limit

	return limit, offset, page
}

func (h *CollectorHandler) isValidJobStatus(status string, validStatuses []string) bool {
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func (h *CollectorHandler) RegisterRoutes(app *fiber.App) {

	collectors := app.Group("/api/v1/collectors")

	collectors.Post("/", h.CreateCollector)
	collectors.Get("/:id", h.GetCollectorByID)
	collectors.Put("/:id", h.UpdateCollector)
	collectors.Delete("/:id", h.DeleteCollector)

	collectors.Get("/", h.ListCollectors)
	collectors.Get("/active", h.GetActiveCollectors)
	collectors.Get("/user/:userID", h.GetCollectorByUserID)
	collectors.Get("/address/:addressID", h.GetCollectorsByAddress)
	collectors.Get("/trash-category/:trashCategoryID", h.GetCollectorsByTrashCategory)

	collectors.Patch("/:id/job-status", h.UpdateJobStatus)
	collectors.Patch("/:id/rating", h.UpdateRating)
	collectors.Put("/:id/available-trash", h.UpdateAvailableTrash)
}
