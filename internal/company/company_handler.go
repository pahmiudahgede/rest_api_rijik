package company

import (
	"context"
	"log"
	"rijig/middleware"
	"rijig/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CompanyProfileHandler struct {
	service CompanyProfileService
}

func NewCompanyProfileHandler(service CompanyProfileService) *CompanyProfileHandler {
	return &CompanyProfileHandler{
		service: service,
	}
}

func (h *CompanyProfileHandler) CreateCompanyProfile(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		return utils.Unauthorized(c, "unauthorized access")
	}

	var req RequestCompanyProfileDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	if errors, valid := req.ValidateCompanyProfileInput(); !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "validation failed", errors)
	}

	companyLogo, err := c.FormFile("company_logo")
	if err != nil {
		log.Printf("Error getting company logo: %v", err)
		return utils.BadRequest(c, "company logo is required")
	}

	res, err := h.service.CreateCompanyProfile(c.Context(), claims.UserID, claims.DeviceID, &req, companyLogo)
	if err != nil {
		log.Printf("Error creating identity card: %v", err)
		if strings.Contains(err.Error(), "invalid file type") {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create company logo")
	}

	return utils.SuccessWithData(c, "company profile created successfully", res)
}

func (h *CompanyProfileHandler) GetCompanyProfileByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "User not authenticated")
	}

	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "id is required")
	}

	res, err := h.service.GetCompanyProfileByID(context.Background(), id)
	if err != nil {
		return utils.NotFound(c, err.Error())
	}

	return utils.SuccessWithData(c, "company profile retrieved", res)
}

func (h *CompanyProfileHandler) GetCompanyProfilesByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "User not authenticated")
	}

	res, err := h.service.GetCompanyProfilesByUserID(context.Background(), userID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "company profiles retrieved", res)
}

func (h *CompanyProfileHandler) UpdateCompanyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "User not authenticated")
	}

	var req RequestCompanyProfileDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	if errors, valid := req.ValidateCompanyProfileInput(); !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "validation failed", errors)
	}

	res, err := h.service.UpdateCompanyProfile(context.Background(), userID, &req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "company profile updated", res)
}

func (h *CompanyProfileHandler) DeleteCompanyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "User not authenticated")
	}

	err := h.service.DeleteCompanyProfile(context.Background(), userID)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "company profile deleted")
}
