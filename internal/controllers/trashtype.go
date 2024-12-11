package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetTrashCategories(c *fiber.Ctx) error {
	trashCategories, err := services.GetTrashCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch trash categories",
			nil,
		))
	}

	var response []dto.TrashCategoryResponse
	for _, category := range trashCategories {

		response = append(response, dto.NewTrashCategoryResponse(
			category.ID,
			category.Name,
			utils.FormatDateToIndonesianFormat(category.CreatedAt),
			utils.FormatDateToIndonesianFormat(category.UpdatedAt),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash categories fetched successfully",
		response,
	))
}

func GetTrashCategoryDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	trashCategoryDetail, err := services.GetTrashCategoryDetail(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch trash category detail",
			nil,
		))
	}

	detailResponse := dto.NewTrashCategoryResponse(
		trashCategoryDetail.ID,
		trashCategoryDetail.Name,
		utils.FormatDateToIndonesianFormat(trashCategoryDetail.CreatedAt),
		utils.FormatDateToIndonesianFormat(trashCategoryDetail.UpdatedAt),
	)

	var detailResponseList []dto.TrashDetailResponse
	if trashCategoryDetail.Details != nil {
		for _, detail := range trashCategoryDetail.Details {
			detailResponseList = append(detailResponseList, dto.NewTrashDetailResponse(
				detail.ID,
				detail.Description,
				detail.Price,
				utils.FormatDateToIndonesianFormat(detail.CreatedAt),
				utils.FormatDateToIndonesianFormat(detail.UpdatedAt),
			))
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash category detail fetched successfully",
		struct {
			Category dto.TrashCategoryResponse `json:"category"`
			Details  []dto.TrashDetailResponse `json:"details,omitempty"`
		}{
			Category: detailResponse,
			Details:  detailResponseList,
		},
	))
}

func CreateTrashCategory(c *fiber.Ctx) error {
	var categoryInput dto.TrashCategoryDTO

	if err := c.BodyParser(&categoryInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := categoryInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	newCategory, err := services.CreateTrashCategory(categoryInput.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create trash category",
			nil,
		))
	}

	categoryResponse := map[string]interface{}{
		"id":        newCategory.ID,
		"name":      newCategory.Name,
		"createdAt": newCategory.CreatedAt,
		"updatedAt": newCategory.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Trash category created successfully",
		categoryResponse,
	))
}

func CreateTrashDetail(c *fiber.Ctx) error {
	var detailInput dto.TrashDetailDTO

	if err := c.BodyParser(&detailInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := detailInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	newDetail, err := services.CreateTrashDetail(detailInput.CategoryID, detailInput.Description, detailInput.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create trash detail",
			nil,
		))
	}

	detailResponse := map[string]interface{}{
		"id":          newDetail.ID,
		"description": newDetail.Description,
		"price":       newDetail.Price,
		"createdAt":   newDetail.CreatedAt,
		"updatedAt":   newDetail.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Trash detail created successfully",
		detailResponse,
	))
}

func UpdateTrashCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var categoryInput dto.UpdateTrashCategoryDTO
	if err := c.BodyParser(&categoryInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := categoryInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	updatedCategory, err := services.UpdateTrashCategory(id, categoryInput.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update trash category",
			nil,
		))
	}

	response := dto.NewTrashCategoryResponse(
		updatedCategory.ID,
		updatedCategory.Name,
		utils.FormatDateToIndonesianFormat(updatedCategory.CreatedAt),
		utils.FormatDateToIndonesianFormat(updatedCategory.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash category updated successfully",
		response,
	))
}

func UpdateTrashDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	var detailInput dto.UpdateTrashDetailDTO
	if err := c.BodyParser(&detailInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input data",
			nil,
		))
	}

	if err := detailInput.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation failed: "+err.Error(),
			nil,
		))
	}

	updatedDetail, err := services.UpdateTrashDetail(id, detailInput.Description, detailInput.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update trash detail",
			nil,
		))
	}

	response := dto.NewTrashDetailResponse(
		updatedDetail.ID,
		updatedDetail.Description,
		updatedDetail.Price,
		utils.FormatDateToIndonesianFormat(updatedDetail.CreatedAt),
		utils.FormatDateToIndonesianFormat(updatedDetail.UpdatedAt),
	)

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash detail updated successfully",
		response,
	))
}

func DeleteTrashCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	err := services.DeleteTrashCategory(id)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete trash category",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash category deleted successfully",
		nil,
	))
}

func DeleteTrashDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	err := services.DeleteTrashDetail(id)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete trash detail",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Trash detail deleted successfully",
		nil,
	))
}
