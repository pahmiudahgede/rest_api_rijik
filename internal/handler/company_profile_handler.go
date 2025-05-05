package handler

import (
	"fmt"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CompanyProfileHandler struct {
	companyProfileService services.CompanyProfileService
}

func NewCompanyProfileHandler(service services.CompanyProfileService) *CompanyProfileHandler {
	return &CompanyProfileHandler{
		companyProfileService: service,
	}
}

func (h *CompanyProfileHandler) CreateCompanyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	var requestDTO dto.RequestCompanyProfileDTO
	if err := c.BodyParser(&requestDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid input data"}})
	}

	companyProfileResponse, err := h.companyProfileService.CreateCompanyProfile(userID, &requestDTO)
	if err != nil {
		return utils.ErrorResponse(c, fmt.Sprintf("Failed to create company profile: %v", err))
	}

	return utils.SuccessResponse(c, companyProfileResponse, "Company profile created successfully")
}

func (h *CompanyProfileHandler) GetCompanyProfileByID(c *fiber.Ctx) error {
	id := c.Params("company_id")

	companyProfileResponse, err := h.companyProfileService.GetCompanyProfileByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fmt.Sprintf("Failed to fetch company profile: %v", err))
	}

	return utils.SuccessResponse(c, companyProfileResponse, "Company profile fetched successfully")
}

func (h *CompanyProfileHandler) GetCompanyProfilesByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	companyProfilesResponse, err := h.companyProfileService.GetCompanyProfilesByUserID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fmt.Sprintf("Failed to fetch company profiles: %v", err))
	}

	return utils.NonPaginatedResponse(c, companyProfilesResponse, len(companyProfilesResponse), "Company profiles fetched successfully")
}

func (h *CompanyProfileHandler) UpdateCompanyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}

	id := c.Params("company_id")

	var requestDTO dto.RequestCompanyProfileDTO
	if err := c.BodyParser(&requestDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid input data"}})
	}

	companyProfileResponse, err := h.companyProfileService.UpdateCompanyProfile(id, &requestDTO)
	if err != nil {
		return utils.ErrorResponse(c, fmt.Sprintf("Failed to update company profile: %v", err))
	}

	return utils.SuccessResponse(c, companyProfileResponse, "Company profile updated successfully")
}

func (h *CompanyProfileHandler) DeleteCompanyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.GenericResponse(c, fiber.StatusUnauthorized, "Unauthorized: User session not found")
	}
	id := c.Params("company_id")

	err := h.companyProfileService.DeleteCompanyProfile(id)
	if err != nil {
		return utils.ErrorResponse(c, fmt.Sprintf("Failed to delete company profile: %v", err))
	}

	return utils.SuccessResponse(c, nil, "Company profile deleted successfully")
}
