package about

import (
	"fmt"
	"log"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type AboutHandler struct {
	AboutService services.AboutService
}

func NewAboutHandler(aboutService services.AboutService) *AboutHandler {
	return &AboutHandler{
		AboutService: aboutService,
	}
}

func (h *AboutHandler) CreateAbout(c *fiber.Ctx) error {
	var request dto.RequestAboutDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateAbout()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	aboutCoverImage, err := c.FormFile("cover_image")
	if err != nil {
		return utils.BadRequest(c, "Cover image is required")
	}

	response, err := h.AboutService.CreateAbout(request, aboutCoverImage)
	if err != nil {
		log.Printf("Error creating About: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to create About: %v", err))
	}

	return utils.CreateSuccessWithData(c, "Successfully created About", response)
}

func (h *AboutHandler) UpdateAbout(c *fiber.Ctx) error {
	id := c.Params("id")

	var request dto.RequestAboutDTO
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.BadRequest(c, "Invalid input data")
	}

	aboutCoverImage, err := c.FormFile("cover_image")
	if err != nil {
		log.Printf("Error retrieving cover image about from request: %v", err)
		return utils.BadRequest(c, "cover_image is required")
	}

	response, err := h.AboutService.UpdateAbout(id, request, aboutCoverImage)
	if err != nil {
		log.Printf("Error updating About: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to update About: %v", err))
	}

	return utils.SuccessWithData(c, "Successfully updated About", response)
}

func (h *AboutHandler) GetAllAbout(c *fiber.Ctx) error {
	response, err := h.AboutService.GetAllAbout()
	if err != nil {
		log.Printf("Error fetching all About: %v", err)
		return utils.InternalServerError(c, "Failed to fetch About list")
	}

	return utils.SuccessWithPagination(c, "Successfully fetched About list", response, 1, len(response))
}

func (h *AboutHandler) GetAboutByID(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := h.AboutService.GetAboutByID(id)
	if err != nil {
		log.Printf("Error fetching About by ID: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to fetch About by ID: %v", err))
	}

	return utils.SuccessWithData(c, "Successfully fetched About", response)
}

func (h *AboutHandler) GetAboutDetailById(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := h.AboutService.GetAboutDetailById(id)
	if err != nil {
		log.Printf("Error fetching About detail by ID: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to fetch About by ID: %v", err))
	}

	return utils.SuccessWithData(c, "Successfully fetched About", response)
}

func (h *AboutHandler) DeleteAbout(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.AboutService.DeleteAbout(id); err != nil {
		log.Printf("Error deleting About: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to delete About: %v", err))
	}

	return utils.Success(c, "Successfully deleted About")
}

func (h *AboutHandler) CreateAboutDetail(c *fiber.Ctx) error {
	var request dto.RequestAboutDetailDTO
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.BadRequest(c, "Invalid input data")
	}

	errors, valid := request.ValidateAboutDetail()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	aboutDetailImage, err := c.FormFile("image_detail")
	if err != nil {
		log.Printf("Error retrieving image detail from request: %v", err)
		return utils.BadRequest(c, "image_detail is required")
	}

	response, err := h.AboutService.CreateAboutDetail(request, aboutDetailImage)
	if err != nil {
		log.Printf("Error creating AboutDetail: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to create AboutDetail: %v", err))
	}

	return utils.CreateSuccessWithData(c, "Successfully created AboutDetail", response)
}

func (h *AboutHandler) UpdateAboutDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	var request dto.RequestAboutDetailDTO
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.BadRequest(c, "Invalid input data")
	}

	aboutDetailImage, err := c.FormFile("image_detail")
	if err != nil {
		log.Printf("Error retrieving image detail from request: %v", err)
		return utils.BadRequest(c, "image_detail is required")
	}

	response, err := h.AboutService.UpdateAboutDetail(id, request, aboutDetailImage)
	if err != nil {
		log.Printf("Error updating AboutDetail: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to update AboutDetail: %v", err))
	}

	return utils.SuccessWithData(c, "Successfully updated AboutDetail", response)
}

func (h *AboutHandler) DeleteAboutDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.AboutService.DeleteAboutDetail(id); err != nil {
		log.Printf("Error deleting AboutDetail: %v", err)
		return utils.InternalServerError(c, fmt.Sprintf("Failed to delete AboutDetail: %v", err))
	}

	return utils.Success(c, "Successfully deleted AboutDetail")
}
